package gemini

import (
	"context"
	"os"
	"testing"
)

func TestGenerateImagen3Image(t *testing.T) {
	if os.Getenv("GEMINI_API_KEY") == "" {
		t.Skip("GEMINI_API_KEY not set")
	}

	svc := NewGeminiService()
	img, err := svc.GenerateImagen3Image(context.Background(), "A red apple on a table")
	if err != nil {
		t.Fatalf("GenerateImagen3Image returned error: %v", err)
	}
	if len(img) == 0 {
		t.Fatalf("GenerateImagen3Image returned empty image")
	}
}

func TestGenerateFlash2Image(t *testing.T) {
	if os.Getenv("GEMINI_API_KEY") == "" {
		t.Skip("GEMINI_API_KEY not set")
	}

	svc := NewGeminiService()
	img, err := svc.GenerateFlash2Image(context.Background(), "A blue cube")
	if err != nil {
		t.Fatalf("GenerateFlash2Image returned error: %v", err)
	}
	if len(img) == 0 {
		t.Fatalf("GenerateFlash2Image returned empty image")
	}
}
