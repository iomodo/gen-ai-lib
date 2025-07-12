package openai

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"

	goopenai "github.com/sashabaranov/go-openai"
)

const (
	GPTImage1 = goopenai.CreateImageModelGptImage1
	DallE3    = goopenai.CreateImageModelDallE3

	GPT4OMini = goopenai.GPT4oMini
	GPT41Mini = goopenai.GPT4Dot1Mini
)

// OpenAIService provides helpers around the go-openai client.
type OpenAIService interface {
	GenerateGPTImage1(ctx context.Context, prompt string) ([]byte, error)
	GenerateGPTImage1WithImage(ctx context.Context, prompt, imageURL string) ([]byte, error)
	GenerateDallEImage(ctx context.Context, prompt, size string) ([]byte, error)
	GenerateResponseFromContent(ctx context.Context, content string) (string, error)
	SanitizePrompt(ctx context.Context, prompt string) (string, error)
	Moderation(ctx context.Context, text string) (bool, error)
}

type service struct {
	client *goopenai.Client
}

// NewService returns an OpenAIService using the OPENAI_API_KEY environment variable.
func NewService() OpenAIService {
	apiKey := os.Getenv("OPENAI_API_KEY")
	cli := goopenai.NewClient(apiKey)
	return &service{client: cli}
}

func (s *service) GenerateGPTImage1(ctx context.Context, prompt string) ([]byte, error) {
	req := goopenai.ImageRequest{
		Model:          GPTImage1,
		Prompt:         prompt,
		Moderation:     goopenai.CreateImageModerationLow,
		OutputFormat:   goopenai.CreateImageOutputFormatWEBP,
		Quality:        goopenai.CreateImageQualityHigh,
		Size:           goopenai.CreateImageSize1024x1024,
		ResponseFormat: goopenai.CreateImageResponseFormatB64JSON,
		N:              1,
	}
	resp, err := s.client.CreateImage(ctx, req)
	if err != nil {
		return nil, err
	}
	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("empty response")
	}
	buf, err := base64.StdEncoding.DecodeString(resp.Data[0].B64JSON)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (s *service) GenerateGPTImage1WithImage(ctx context.Context, prompt, imageURL string) ([]byte, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" || contentType == "application/octet-stream" {
		switch {
		case hasSuffix(imageURL, ".jpeg"), hasSuffix(imageURL, ".jpg"):
			contentType = "image/jpeg"
		case hasSuffix(imageURL, ".png"):
			contentType = "image/png"
		case hasSuffix(imageURL, ".webp"):
			contentType = "image/webp"
		default:
			contentType = "image/png"
		}
	}

	reader := goopenai.WrapReader(bytes.NewReader(data), "image"+extFromContentType(contentType), contentType)
	req := goopenai.ImageEditRequest{
		Model:          GPTImage1,
		Prompt:         prompt,
		Image:          reader,
		Quality:        goopenai.CreateImageQualityHigh,
		Size:           goopenai.CreateImageSize1024x1024,
		ResponseFormat: goopenai.CreateImageResponseFormatB64JSON,
		N:              1,
	}
	editResp, err := s.client.CreateEditImage(ctx, req)
	if err != nil {
		return nil, err
	}
	if len(editResp.Data) == 0 {
		return nil, fmt.Errorf("empty response")
	}
	buf, err := base64.StdEncoding.DecodeString(editResp.Data[0].B64JSON)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (s *service) GenerateDallEImage(ctx context.Context, prompt, size string) ([]byte, error) {
	req := goopenai.ImageRequest{
		Model:          DallE3,
		Prompt:         prompt,
		Quality:        goopenai.CreateImageQualityHD,
		Size:           size,
		ResponseFormat: goopenai.CreateImageResponseFormatB64JSON,
		N:              1,
	}
	resp, err := s.client.CreateImage(ctx, req)
	if err != nil {
		return nil, err
	}
	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("empty response")
	}
	buf, err := base64.StdEncoding.DecodeString(resp.Data[0].B64JSON)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (s *service) GenerateResponseFromContent(ctx context.Context, content string) (string, error) {
	req := goopenai.ChatCompletionRequest{
		Model:    GPT4OMini,
		Messages: []goopenai.ChatCompletionMessage{{Role: goopenai.ChatMessageRoleUser, Content: content}},
	}
	resp, err := s.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("empty response")
	}
	return resp.Choices[0].Message.Content, nil
}

func (s *service) SanitizePrompt(ctx context.Context, prompt string) (string, error) {
	text := fmt.Sprintf(`Please rewrite the following image generation prompt to be safe and appropriate while maintaining the core creative intent. Follow these rules:
1. Replace any specific brand names, IP, or copyrighted content with generic descriptions
2. Remove or replace any potentially offensive, sexual, or violent content
3. Keep the artistic style and main subject matter intact
4. Make the description more general while preserving the creative vision
5. Ensure the prompt follows content policy guidelines
6. Keep the response concise and focused on visual elements

Original prompt: %s`, prompt)

	req := goopenai.ChatCompletionRequest{
		Model: GPT41Mini,
		Messages: []goopenai.ChatCompletionMessage{
			{Role: goopenai.ChatMessageRoleSystem, Content: "You are a prompt sanitization assistant. Your task is to rewrite image generation prompts to be safe and appropriate while maintaining the creative intent. Always respond with just the sanitized prompt, no explanations or additional text."},
			{Role: goopenai.ChatMessageRoleUser, Content: text},
		},
	}
	resp, err := s.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("empty response")
	}
	return resp.Choices[0].Message.Content, nil
}

func (s *service) Moderation(ctx context.Context, text string) (bool, error) {
	resp, err := s.client.Moderations(ctx, goopenai.ModerationRequest{
		Model: goopenai.ModerationOmniLatest,
		Input: text,
	})
	if err != nil {
		return false, err
	}
	if len(resp.Results) == 0 {
		return false, nil
	}
	return resp.Results[0].Flagged, nil
}

func hasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

func extFromContentType(ct string) string {
	switch ct {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	}
	return ""
}
