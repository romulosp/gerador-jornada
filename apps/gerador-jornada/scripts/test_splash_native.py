"""Verifica imagem real em inicializações consecutivas; exige Windows gráfico e Pillow no desenvolvimento."""

import argparse
import ctypes
from ctypes import wintypes
import json
import os
from pathlib import Path
import subprocess
import time


def capture_window(process_id, target):
    """Captura somente a janela visível do processo verificado e recusa vídeo preto/minúsculo."""
    from PIL import ImageGrab

    # As assinaturas mantêm HWND de 64 bits e limitam a busca ao processo iniciado pelo teste.
    user = ctypes.WinDLL("user32")
    user.GetWindowThreadProcessId.argtypes = [wintypes.HWND, ctypes.POINTER(wintypes.DWORD)]
    user.GetWindowRect.argtypes = [wintypes.HWND, ctypes.POINTER(wintypes.RECT)]
    user.GetClientRect.argtypes = [wintypes.HWND, ctypes.POINTER(wintypes.RECT)]
    user.ClientToScreen.argtypes = [wintypes.HWND, ctypes.POINTER(wintypes.POINT)]
    user.GetWindowLongW.argtypes = [wintypes.HWND, ctypes.c_int]
    user.GetClassNameW.argtypes = [wintypes.HWND, wintypes.LPWSTR, ctypes.c_int]
    user.IsWindowVisible.argtypes = [wintypes.HWND]
    callback_type = ctypes.WINFUNCTYPE(wintypes.BOOL, wintypes.HWND, wintypes.LPARAM)
    windows = []

    @callback_type
    def visit(handle, _context):
        """Seleciona apenas os HWND visíveis pertencentes ao app sob teste."""
        owner = wintypes.DWORD()
        user.GetWindowThreadProcessId(handle, ctypes.byref(owner))
        if owner.value == process_id and user.IsWindowVisible(handle):
            windows.append(handle)
        return True

    user.EnumWindows(visit, 0)
    if not windows:
        raise AssertionError("Janela de splash ausente")
    # A amostra tem limiar específico da arte v2; eventos PLAY sozinhos não provam exibição.
    rectangle = wintypes.RECT()
    user.GetWindowRect(windows[0], ctypes.byref(rectangle))
    # A apresentação deve ter moldura normal e ocupar menos que o monitor, sem maximizar/fullscreen.
    class MonitorInfo(ctypes.Structure):
        _fields_ = [("size", wintypes.DWORD), ("monitor", wintypes.RECT), ("work", wintypes.RECT), ("flags", wintypes.DWORD)]

    user.MonitorFromWindow.argtypes = [wintypes.HWND, wintypes.DWORD]
    user.MonitorFromWindow.restype = wintypes.HANDLE
    user.GetMonitorInfoW.argtypes = [wintypes.HANDLE, ctypes.POINTER(MonitorInfo)]
    monitor = MonitorInfo()
    monitor.size = ctypes.sizeof(monitor)
    assert user.GetMonitorInfoW(user.MonitorFromWindow(windows[0], 2), ctypes.byref(monitor)), "Monitor da splash não encontrado"
    assert any(abs(getattr(rectangle, side)-getattr(monitor.monitor, side)) > 1 for side in ("left", "top", "right", "bottom")), "Splash entrou em tela cheia"
    style = user.GetWindowLongW(windows[0], -16)
    assert style & 0x00C00000 == 0x00C00000, "Barra de título da splash ausente"
    assert not style & 0x01000000, "Splash abriu maximizada"
    # A superfície deve ocupar a maior área proporcional da janela, inclusive em monitores maiores que o MP4.
    children = []

    @callback_type
    def visit_child(handle, _context):
        """Inspeciona a geometria da própria superfície de vídeo, sem inferir pelo evento PLAY."""
        name = ctypes.create_unicode_buffer(128)
        user.GetClassNameW(handle, name, len(name))
        if name.value == "JornadaSplashVideoSurface":
            bounds = wintypes.RECT()
            user.GetWindowRect(handle, ctypes.byref(bounds))
            children.append(bounds)
        return True

    user.EnumChildWindows.argtypes = [wintypes.HWND, callback_type, wintypes.LPARAM]
    user.EnumChildWindows(windows[0], visit_child, 0)
    assert len(children) == 1, "Superfície de vídeo ausente"
    # O vídeo se ajusta ao cliente, sem contar bordas e barra de título como área de reprodução.
    client = wintypes.RECT()
    user.GetClientRect(windows[0], ctypes.byref(client))
    origin = wintypes.POINT()
    user.ClientToScreen(windows[0], ctypes.byref(origin))
    width, height = client.right-client.left, client.bottom-client.top
    scale = min(width/1280, height/720)
    bounds = children[0]
    assert abs(bounds.right-bounds.left-int(1280*scale)) <= 1, "Vídeo não ocupa a maior largura proporcional"
    assert abs(bounds.bottom-bounds.top-int(720*scale)) <= 1, "Vídeo não ocupa a maior altura proporcional"
    assert bounds.left >= origin.x and bounds.top >= origin.y, "Vídeo ultrapassa o início do cliente"
    assert bounds.right <= origin.x+width and bounds.bottom <= origin.y+height, "Vídeo ultrapassa o fim do cliente"
    shot = ImageGrab.grab(bbox=(rectangle.left, rectangle.top, rectangle.right, rectangle.bottom))
    shot.save(target)
    sample = shot.convert("RGB").resize((200, 100))
    colorful = sum(1 for red, green, blue in sample.get_flattened_data() if max(red, green, blue)-min(red, green, blue)>45)
    if colorful <= 1500:
        raise AssertionError(f"Vídeo ausente/minúsculo: {colorful} pixels coloridos")
    return colorful


def startup_records(process_id):
    """Lê somente a sessão técnica do processo de teste, sem nomes/conteúdo dos projetos."""
    directory = Path(os.environ["LOCALAPPDATA"]) / "GeradorJornada" / "logs"
    for source in sorted(directory.glob("*.jsonl"), key=lambda path: path.stat().st_mtime, reverse=True):
        records = [json.loads(line) for line in source.read_text(encoding="utf-8").splitlines() if line.strip()]
        if records and records[0].get("pid") == process_id:
            return records
    raise AssertionError("Sessão técnica do teste não encontrada")


def main():
    """Executa repetições sequenciais com áudio, capturas em dois instantes e verificação da barreira."""
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--runs", type=int, default=4)
    arguments = parser.parse_args()
    if os.name != "nt":
        parser.error("Este smoke exige sessão gráfica Windows")
    root = Path(__file__).resolve().parents[1]
    evidence = root / "build" / "036-v11-repeated"
    evidence.mkdir(parents=True, exist_ok=True)
    environment = os.environ.copy()
    environment["JORNADA_CONFIG_DIR"] = str(evidence / "config")
    # Este teste verifica o renderer repetido; a política de primeiro uso tem smoke próprio.
    config = Path(environment["JORNADA_CONFIG_DIR"])
    config.mkdir(parents=True, exist_ok=True)
    (config / "startup-preferences.json").write_text(json.dumps({"presented": True, "show_on_startup": True}), encoding="utf-8")
    executable = root / "dist" / "windows-amd64" / "JornadaStudio.exe"
    for index in range(1, arguments.runs+1):
        with (evidence/f"process-{index}.log").open("w", encoding="utf-8") as output:
            process = subprocess.Popen([str(executable), "-smoke-screenshot", str(evidence/f"launcher-{index}.png")], env=environment, stdout=output, stderr=subprocess.STDOUT)
            try:
                time.sleep(2)
                early = capture_window(process.pid, evidence/f"video-{index}-early.png")
                time.sleep(6)
                late = capture_window(process.pid, evidence/f"video-{index}-late.png")
                process.wait(timeout=20)
                assert process.returncode == 0, "Processo encerrou com falha"
                records = startup_records(process.pid)
                events = [record["event"] for record in records]
                assert events.index("startup_resources_ready") < events.index("splash_video_ended") < events.index("startup_ready")
                assert next(record for record in records if record["event"]=="splash_video_started")["audio_track"]
                (evidence/f"startup-{index}.jsonl").write_text("\n".join(json.dumps(record,ensure_ascii=False) for record in records)+"\n",encoding="utf-8")
                print(f"PASS {index}: imagem inicial={early}, final={late}; áudio e término confirmados",flush=True)
            finally:
                if process.poll() is None:
                    process.terminate()
                    process.wait(timeout=5)


if __name__ == "__main__":
    main()
