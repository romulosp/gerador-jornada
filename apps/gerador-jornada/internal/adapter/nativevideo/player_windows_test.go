package nativevideo

import "testing"

// TestEntireFrameFits verifica a própria função nativa, incluindo telas pequenas, grandes e proporções distintas.
func TestEntireFrameFits(t *testing.T) {
	if size := DisplaySize(1280, 720, 16, 9); size != [2]int{1280, 720} {
		t.Fatalf("razão foi tratada como resolução: %v", size)
	}
	if size := DisplaySize(720, 576, 4, 3); size != [2]int{768, 576} {
		t.Fatalf("correção de pixels não quadrados incorreta: %v", size)
	}
	// Cada tamanho representa pixels físicos já obtidos do HWND, independentemente do DPI lógico do Fyne.
	for _, size := range [][2]int{{1366, 768}, {1920, 1080}, {800, 600}, {720, 1280}, {3840, 2160}, {1, 1}} {
		rect := FitRect(size[0], size[1], 1280, 720)
		if rect[0] < 0 || rect[1] < 0 || rect[0]+rect[2] > size[0] || rect[1]+rect[3] > size[1] {
			t.Fatalf("quadro cortado em %v: %v", size, rect)
		}
	}
	if rect := FitRect(1920, 1080, 1280, 720); rect != [4]int{0, 0, 1920, 1080} {
		t.Fatalf("vídeo não ocupa a tela inteira: %v", rect)
	}
	if rect := FitRect(3840, 2160, 1280, 720); rect != [4]int{0, 0, 3840, 2160} {
		t.Fatalf("vídeo não ocupa a tela 4K: %v", rect)
	}
	if rect := FitRect(800, 600, 1280, 720); rect != [4]int{0, 75, 800, 450} {
		t.Fatalf("redução proporcional incorreta: %v", rect)
	}
}
