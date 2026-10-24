#define WINVER 0x0601
#define _WIN32_WINNT 0x0601
#include <windows.h>
#include <mfplay.h>
#include <atomic>
#include <new>
#include "player_windows.h"

// videoStartMessage serializa a configuração/reprodução na thread que criou o HWND.
static constexpr UINT videoStartMessage = WM_APP + 78;
// videoWindowProcedure atende pintura e tamanho conforme o contrato MFPlay, sem controles STATIC.
static LRESULT CALLBACK videoWindowProcedure(HWND hwnd, UINT message, WPARAM command, LPARAM detail);

// Video recebe eventos MFPlay em threads nativas; as leituras Go usam somente campos atômicos.
class Video final : public IMFPMediaPlayerCallback {
public:
    // refs mantém o callback vivo até o player liberar sua referência COM.
    std::atomic<ULONG> refs{1};
    std::atomic<int> started{0}, ended{0}, audio{0};
    std::atomic<LONG> error{0};
    IMFPMediaPlayer *player = nullptr;
    HWND child = nullptr;
    bool comOwned = false;
    // parentClipped informa se a superfície adicionou proteção de pintura ao HWND pai.
    bool parentClipped = false;

    // QueryInterface oferece apenas IUnknown e o callback documentado do MFPlay.
    HRESULT STDMETHODCALLTYPE QueryInterface(REFIID iid, void **result) override {
        if (!result) return E_POINTER;
        *result = nullptr;
        if (IsEqualIID(iid, IID_IUnknown) || IsEqualIID(iid, __uuidof(IMFPMediaPlayerCallback))) {
            *result = static_cast<IMFPMediaPlayerCallback *>(this); AddRef(); return S_OK;
        }
        return E_NOINTERFACE;
    }
    // AddRef e Release seguem a propriedade compartilhada do callback COM.
    ULONG STDMETHODCALLTYPE AddRef() override { return ++refs; }
    ULONG STDMETHODCALLTYPE Release() override {
        ULONG remaining = --refs;
        if (!remaining) delete this;
        return remaining;
    }
    // OnMediaPlayerEvent distingue término real, início e erro; tempo fixo nunca libera a splash.
    void STDMETHODCALLTYPE OnMediaPlayerEvent(MFP_EVENT_HEADER *event) override {
        if (FAILED(event->hrEvent)) { error.store(event->hrEvent); return; }
        if (event->eEventType == MFP_EVENT_TYPE_PLAY) started.store(1);
        if (event->eEventType == MFP_EVENT_TYPE_PLAYBACK_ENDED) ended.store(1);
        if (event->eEventType == MFP_EVENT_TYPE_MEDIAITEM_SET) {
            BOOL hasAudio = FALSE, selected = FALSE;
            auto item = reinterpret_cast<MFP_MEDIAITEM_SET_EVENT *>(event)->pMediaItem;
            if (item && SUCCEEDED(item->HasAudio(&hasAudio, &selected))) audio.store(hasAudio && selected);
            // O callback livre não toca o HWND ou o renderer; a thread de UI inicia depois da mídia pronta.
            if (!PostMessageW(child,videoStartMessage,0,0)) error.store(HRESULT_FROM_WIN32(GetLastError()));
        }
    }
};

// videoWindowProcedure repinta o último quadro e configura a mídia somente após sua preparação.
static LRESULT CALLBACK videoWindowProcedure(HWND hwnd, UINT message, WPARAM command, LPARAM detail) {
    auto video = reinterpret_cast<Video *>(GetWindowLongPtrW(hwnd,GWLP_USERDATA));
    if (message==WM_NCCREATE) {
        auto creation = reinterpret_cast<CREATESTRUCTW *>(detail);
        video = static_cast<Video *>(creation->lpCreateParams);
        SetWindowLongPtrW(hwnd,GWLP_USERDATA,reinterpret_cast<LONG_PTR>(video));
    }
    if (message==WM_ERASEBKGND) return 1;
    if (message==WM_PAINT) {
        PAINTSTRUCT paint = {};
        BeginPaint(hwnd,&paint);
        if (video && video->player && video->started.load()) video->player->UpdateVideo();
        else FillRect(paint.hdc,&paint.rcPaint,static_cast<HBRUSH>(GetStockObject(BLACK_BRUSH)));
        EndPaint(hwnd,&paint);
        return 0;
    }
    if (message==videoStartMessage && video && video->player) {
        // A janela e a mídia já existem: aplicar origem integral, letterbox e áudio antes do primeiro Play.
        MFVideoNormalizedRect full = {0.0f,0.0f,1.0f,1.0f};
        HRESULT hr = video->player->SetVideoSourceRect(&full);
        if (SUCCEEDED(hr)) hr = video->player->SetAspectRatioMode(MFVideoARMode_PreservePicture);
        if (SUCCEEDED(hr)) hr = video->player->SetMute(FALSE);
        if (SUCCEEDED(hr)) hr = video->player->SetVolume(1.0f);
        if (SUCCEEDED(hr)) {
            jornada_video_resize(video);
            hr = video->player->Play();
        }
        if (FAILED(hr)) video->error.store(hr);
        InvalidateRect(hwnd,nullptr,FALSE);
        return 0;
    }
    if (message==WM_SIZE && video && video->player && video->started.load()) video->player->UpdateVideo();
    if (message==WM_NCDESTROY) SetWindowLongPtrW(hwnd,GWLP_USERDATA,0);
    return DefWindowProcW(hwnd,message,command,detail);
}

// jornada_video_open cria uma superfície filha e reproduz o MP4 original com áudio habilitado.
extern "C" void *jornada_video_open(uintptr_t parent, const uint16_t *path, int32_t *error) {
    HRESULT hr = CoInitializeEx(nullptr, COINIT_APARTMENTTHREADED);
    bool owned = SUCCEEDED(hr);
    if (FAILED(hr) && hr != RPC_E_CHANGED_MODE) { *error = hr; return nullptr; }
    auto video = new (std::nothrow) Video();
    if (!video) { if (owned) CoUninitialize(); *error = E_OUTOFMEMORY; return nullptr; }
    video->comOwned = owned;
    // GetClientRect é a referência física: multiplicar novamente o tamanho Fyne por DPI pode cortar a imagem.
    RECT client = {};
    if (!GetClientRect(reinterpret_cast<HWND>(parent), &client)) { *error = HRESULT_FROM_WIN32(GetLastError()); jornada_video_close(video); return nullptr; }
    int width = client.right-client.left, height = client.bottom-client.top;
    LONG_PTR parentStyle = GetWindowLongPtrW(reinterpret_cast<HWND>(parent),GWL_STYLE);
    video->parentClipped = !(parentStyle & WS_CLIPCHILDREN);
    if (video->parentClipped) SetWindowLongPtrW(reinterpret_cast<HWND>(parent),GWL_STYLE,parentStyle|WS_CLIPCHILDREN);
    // A classe própria atende WM_PAINT/WM_SIZE; uma janela STATIC pode apagar o renderer de forma intermitente.
    WNDCLASSW type = {};
    type.lpfnWndProc = videoWindowProcedure;
    type.hInstance = GetModuleHandleW(nullptr);
    type.hbrBackground = static_cast<HBRUSH>(GetStockObject(BLACK_BRUSH));
    type.lpszClassName = L"JornadaSplashVideoSurface";
    if (!RegisterClassW(&type) && GetLastError()!=ERROR_CLASS_ALREADY_EXISTS) {
        *error=HRESULT_FROM_WIN32(GetLastError()); jornada_video_close(video); return nullptr;
    }
    video->child = CreateWindowExW(0, type.lpszClassName, L"", WS_CHILD | WS_VISIBLE | WS_CLIPSIBLINGS,
        0, 0, width, height, reinterpret_cast<HWND>(parent), nullptr, type.hInstance, video);
    if (!video->child) { *error = HRESULT_FROM_WIN32(GetLastError()); jornada_video_close(video); return nullptr; }
    hr = MFPCreateMediaPlayer(reinterpret_cast<LPCWSTR>(path), FALSE, MFP_OPTION_FREE_THREADED_CALLBACK,
        video, video->child, &video->player);
    if (SUCCEEDED(hr)) hr = video->player->SetMute(FALSE);
    if (SUCCEEDED(hr)) hr = video->player->SetVolume(1.0f);
    if (FAILED(hr)) { *error = hr; jornada_video_close(video); return nullptr; }
    *error = 0;
    return video;
}

// jornada_video_snapshot publica o estado sem depender da thread do callback.
extern "C" void jornada_video_snapshot(void *player, JornadaVideoSnapshot *snapshot) {
    auto video = static_cast<Video *>(player);
    snapshot->started = video->started.load(); snapshot->ended = video->ended.load();
    snapshot->audio = video->audio.load(); snapshot->error = video->error.load();
}

// jornada_video_fit ocupa a maior área de tela possível preservando o quadro completo.
extern "C" void jornada_video_fit(int width, int height, int sourceWidth, int sourceHeight, int *rect) {
    if (width < 1 || height < 1) { for (int i=0; i<4; ++i) rect[i]=0; return; }
    if (sourceWidth < 1 || sourceHeight < 1) { sourceWidth=width; sourceHeight=height; }
    double scale = static_cast<double>(width)/sourceWidth;
    double vertical = static_cast<double>(height)/sourceHeight;
    if (vertical < scale) scale = vertical;
    int fittedWidth = static_cast<int>(sourceWidth*scale);
    int fittedHeight = static_cast<int>(sourceHeight*scale);
    if (fittedWidth < 1) fittedWidth=1;
    if (fittedHeight < 1) fittedHeight=1;
    rect[0]=(width-fittedWidth)/2; rect[1]=(height-fittedHeight)/2;
    rect[2]=fittedWidth; rect[3]=fittedHeight;
}

// jornada_video_resize usa o cliente HWND real e centraliza o quadro completo sem cortar suas bordas.
// jornada_video_display_size recebe a razão retornada pelo MFPlay e conserva a altura nativa.
extern "C" void jornada_video_display_size(int width, int height, int aspectWidth, int aspectHeight, int *size) {
    size[0]=width; size[1]=height;
    if (height>0 && aspectWidth>0 && aspectHeight>0) size[0]=MulDiv(height,aspectWidth,aspectHeight);
}

// jornada_video_resize usa o cliente HWND real e centraliza o quadro completo sem cortar suas bordas.
extern "C" int jornada_video_resize(void *player) {
    auto video = static_cast<Video *>(player);
    RECT client = {};
    if (!GetClientRect(GetParent(video->child), &client)) { video->error.store(HRESULT_FROM_WIN32(GetLastError())); return 0; }
    int width = client.right-client.left, height = client.bottom-client.top;
    SIZE source = {}, aspect = {};
    if (video->player) video->player->GetNativeVideoSize(&source, &aspect);
    // pszARVideo é a razão (por exemplo, 16 e 9), não uma segunda resolução em pixels.
    int display[2] = {};
    jornada_video_display_size(source.cx,source.cy,aspect.cx,aspect.cy,display);
    int rect[4] = {};
    jornada_video_fit(width, height, display[0], display[1], rect);
    // Evitar MoveWindow repetido impede apagar as margens já pintadas pelo canvas.
    RECT current = {};
    GetWindowRect(video->child,&current);
    MapWindowPoints(nullptr,GetParent(video->child),reinterpret_cast<POINT *>(&current),2);
    bool changed = current.left!=rect[0] || current.top!=rect[1] || current.right-current.left!=rect[2] || current.bottom-current.top!=rect[3];
    if (changed) MoveWindow(video->child, rect[0], rect[1], rect[2], rect[3], TRUE);
    // A reprodução também invalida a superfície: WM_PAINT encaminha corretamente ao renderer nativo.
    InvalidateRect(video->child,nullptr,FALSE);
    return changed ? 1 : 0;
}

// jornada_video_close interrompe áudio, libera COM e destrói apenas a superfície da splash.
extern "C" void jornada_video_close(void *player) {
    auto video = static_cast<Video *>(player);
    bool owned = video->comOwned;
    if (video->player) { video->player->Shutdown(); video->player->Release(); video->player = nullptr; }
    if (video->child) {
        HWND parent = GetParent(video->child);
        DestroyWindow(video->child);
        if (video->parentClipped) SetWindowLongPtrW(parent,GWL_STYLE,GetWindowLongPtrW(parent,GWL_STYLE)&~WS_CLIPCHILDREN);
    }
    video->Release();
    if (owned) CoUninitialize();
}

// jornada_window_alpha aplica o fade à janela inteira; alfa 255 restaura opacidade completa.
extern "C" int jornada_window_alpha(uintptr_t window, int alpha) {
    HWND hwnd = reinterpret_cast<HWND>(window);
    LONG_PTR style = GetWindowLongPtrW(hwnd, GWL_EXSTYLE);
    SetWindowLongPtrW(hwnd, GWL_EXSTYLE, style | WS_EX_LAYERED);
    return SetLayeredWindowAttributes(hwnd, 0, static_cast<BYTE>(alpha), LWA_ALPHA) != 0;
}
