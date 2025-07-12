package gemini

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"google.golang.org/genai"
)

const (
	IMAGEN_3_MODEL      = "imagen-3.0-generate-002"
	FLASH_2_MODEL       = "gemini-2.0-flash-exp-image-generation"
	VEO_3_MODEL         = "veo-3.0-video-generation"
	VEO_3_PREVIEW_MODEL = "veo-3.0-generate-preview"
)

type GeminiService interface {
	GenerateImagen3Image(ctx context.Context, prompt string) ([]byte, error)
	GenerateFlash2Image(ctx context.Context, prompt string) ([]byte, error)
	GenerateFlashWithImage(ctx context.Context, prompt, imageURL string) ([]byte, error)
	GenerateVeo3Video(ctx context.Context, prompt string) ([]byte, error)
	GenerateVeo3PreviewVideo(ctx context.Context, prompt string, firstFrame, lastFrame []byte) ([]byte, error)
	GenerateVeo3PreviewVideoFromURLs(ctx context.Context, prompt, firstFrameURL, lastFrameURL string) ([]byte, error)
}

type geminiService struct {
	client *genai.Client
}

func NewGeminiService() GeminiService {
	apiKey := os.Getenv("GEMINI_API_KEY")
	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		log.Fatalf("failed to create genai client: %v", err)
	}
	return &geminiService{
		client: client,
	}
}

func (s *geminiService) GenerateImagen3Image(ctx context.Context, prompt string) ([]byte, error) {
	resp, err := s.client.Models.GenerateImages(ctx, IMAGEN_3_MODEL, prompt, nil)
	if err != nil {
		return nil, err
	}
	if len(resp.GeneratedImages) > 0 && resp.GeneratedImages[0].Image != nil {
		return resp.GeneratedImages[0].Image.ImageBytes, nil
	}
	return nil, nil
}

func (s *geminiService) GenerateFlash2Image(ctx context.Context, prompt string) ([]byte, error) {
	resp, err := s.client.Models.GenerateImages(ctx, FLASH_2_MODEL, prompt, nil)
	if err != nil {
		return nil, err
	}
	if len(resp.GeneratedImages) > 0 && resp.GeneratedImages[0].Image != nil {
		return resp.GeneratedImages[0].Image.ImageBytes, nil
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
	op, err := s.client.Models.GenerateVideos(ctx, VEO_3_MODEL, prompt, nil, nil)
	if err != nil {
		return nil, err
	}
	for !op.Done {
		time.Sleep(2 * time.Second)
		op, err = s.client.Operations.GetVideosOperation(ctx, op, nil)
		if err != nil {
			return nil, err
		}
	}
	if len(op.Response.GeneratedVideos) > 0 {
		data, err := s.client.Files.Download(ctx, genai.NewDownloadURIFromGeneratedVideo(op.Response.GeneratedVideos[0]), nil)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, fmt.Errorf("video generation did not return a result")
}

// GenerateVeo3PreviewVideo creates a video using the veo-3.0-generate-preview model
// by providing the first and last frames along with the prompt.
func (s *geminiService) GenerateVeo3PreviewVideo(ctx context.Context, prompt string, firstFrame, lastFrame []byte) ([]byte, error) {
	start := &genai.Image{ImageBytes: firstFrame, MIMEType: "image/png"}
	cfg := &genai.GenerateVideosConfig{
		LastFrame: &genai.Image{ImageBytes: lastFrame, MIMEType: "image/png"},
	}
	op, err := s.client.Models.GenerateVideos(ctx, VEO_3_PREVIEW_MODEL, prompt, start, cfg)
	if err != nil {
		return nil, err
	}
	for !op.Done {
		time.Sleep(2 * time.Second)
		op, err = s.client.Operations.GetVideosOperation(ctx, op, nil)
		if err != nil {
			return nil, err
		}
	}
	if len(op.Response.GeneratedVideos) > 0 {
		data, err := s.client.Files.Download(ctx, genai.NewDownloadURIFromGeneratedVideo(op.Response.GeneratedVideos[0]), nil)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, fmt.Errorf("video generation did not return a result")
}

// GenerateVeo3PreviewVideoFromURLs downloads the first and last frame images
// from the provided URLs and invokes GenerateVeo3PreviewVideo.
func (s *geminiService) GenerateVeo3PreviewVideoFromURLs(ctx context.Context, prompt, firstFrameURL, lastFrameURL string) ([]byte, error) {
	firstResp, err := http.Get(firstFrameURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download first frame: %w", err)
	}
	defer firstResp.Body.Close()
	firstData, err := io.ReadAll(firstResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read first frame: %w", err)
	}

	lastResp, err := http.Get(lastFrameURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download last frame: %w", err)
	}
	defer lastResp.Body.Close()
	lastData, err := io.ReadAll(lastResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read last frame: %w", err)
	}

	return s.GenerateVeo3PreviewVideo(ctx, prompt, firstData, lastData)
}
