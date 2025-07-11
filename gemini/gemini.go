package gemini

import (
	"context"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

const (
	IMAGEN_3_MODEL = "imagen-3.0-generate-002"
	FLASH_2_MODEL  = "gemini-2.0-flash-exp-image-generation"
	VEO_3_MODEL    = "veo-3.0-video-generation"
)

type GeminiService interface {
	GenerateImagen3Image(ctx context.Context, prompt string) ([]byte, error)
	GenerateFlash2Image(ctx context.Context, prompt string) ([]byte, error)
	GenerateFlashWithImage(ctx context.Context, prompt, imageURL string) ([]byte, error)
	GenerateVeo3Video(ctx context.Context, prompt string) ([]byte, error)
}

type geminiService struct {
	client *genai.Client
}

func NewGeminiService() GeminiService {
	apiKey := os.Getenv("GEMINI_API_KEY")
	client, err := genai.NewClient(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("failed to create genai client: %v", err)
	}
	return &geminiService{
		client: client,
	}
}

func (s *geminiService) GenerateImagen3Image(ctx context.Context, prompt string) ([]byte, error) {
	model := s.client.GenerativeModel(IMAGEN_3_MODEL)
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, err
	}
	for _, part := range resp.Candidates[0].Content.Parts {
		if img, ok := part.(*genai.Blob); ok {
			return img.Data, nil
		}
	}
	return nil, nil
}

func (s *geminiService) GenerateFlash2Image(ctx context.Context, prompt string) ([]byte, error) {
	model := s.client.GenerativeModel(FLASH_2_MODEL)
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, err
	}
	for _, part := range resp.Candidates[0].Content.Parts {
		if img, ok := part.(*genai.Blob); ok {
			return img.Data, nil
		}
	}
	return nil, nil
}

func (s *geminiService) GenerateFlashWithImage(ctx context.Context, prompt, imageURL string) ([]byte, error) {
	// Download image and create a genai.Blob
	// ... (implement image download and blob creation)
	// Then call model.GenerateContent with both prompt and image blob
	return nil, nil // Implement as needed
}

func (s *geminiService) GenerateVeo3Video(ctx context.Context, prompt string) ([]byte, error) {
	model := s.client.GenerativeModel(VEO_3_MODEL)
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, err
	}
	for _, part := range resp.Candidates[0].Content.Parts {
		if video, ok := part.(*genai.Blob); ok {
			return video.Data, nil
		}
	}
	return nil, nil
}
