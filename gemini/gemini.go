package gemini

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/generative-ai-go/genai"
	aiplatform "google.golang.org/api/aiplatform/v1"
	"google.golang.org/api/option"
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

// GenerateVeo3PreviewVideo creates a video using the veo-3.0-generate-preview model
// by providing the first and last frames along with the prompt.
func (s *geminiService) GenerateVeo3PreviewVideo(ctx context.Context, prompt string, firstFrame, lastFrame []byte) ([]byte, error) {
	projectID := os.Getenv("GEMINI_PROJECT_ID")
	svc, err := aiplatform.NewService(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("projects/%s/locations/us-central1/publishers/google/models/%s", projectID, VEO_3_PREVIEW_MODEL)

	req := &aiplatform.GoogleCloudAiplatformV1PredictLongRunningRequest{
		Instances: []interface{}{
			map[string]any{
				"prompt": prompt,
				"first_frame": map[string]any{
					"mimeType":           "image/png",
					"bytesBase64Encoded": base64.StdEncoding.EncodeToString(firstFrame),
				},
				"last_frame": map[string]any{
					"mimeType":           "image/png",
					"bytesBase64Encoded": base64.StdEncoding.EncodeToString(lastFrame),
				},
			},
		},
		Parameters: map[string]any{
			"mime_type": "video/mp4",
		},
	}

	op, err := svc.Projects.Locations.Publishers.Models.PredictLongRunning(endpoint, req).Do()
	if err != nil {
		return nil, err
	}

	fetchReq := &aiplatform.GoogleCloudAiplatformV1FetchPredictOperationRequest{OperationName: op.Name}
	// Poll until the operation is done or context canceled
	for {
		opResp, err := svc.Projects.Locations.Publishers.Models.FetchPredictOperation(endpoint, fetchReq).Do()
		if err != nil {
			return nil, err
		}
		if opResp.Done {
			var data struct {
				Response struct {
					Videos []struct {
						Bytes string `json:"bytesBase64Encoded"`
					} `json:"videos"`
				} `json:"response"`
			}
			if err := json.Unmarshal(opResp.Response, &data); err != nil {
				return nil, err
			}
			if len(data.Response.Videos) > 0 {
				return base64.StdEncoding.DecodeString(data.Response.Videos[0].Bytes)
			}
			break
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(2 * time.Second):
		}
	}
	return nil, fmt.Errorf("video generation did not return a result")
}
