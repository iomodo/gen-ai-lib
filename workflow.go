package genailib

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/iomodo/gen-ai-lib/external/gemini"
	"github.com/iomodo/gen-ai-lib/external/replicate"
	"github.com/pkg/errors"
)

// WorkflowStepFunctionType constants.
const (
	FunctionTypeTextsToText          = "texts_to_text"
	FunctionTypeTextToImage          = "text_to_image"
	FunctionTypeTextAndImageToImage  = "text_and_image_to_image"
	FunctionTypeTextAndImagesToVideo = "text_and_images_to_video"
	FunctionTypeTextAndImageToVideo  = "text_and_image_to_video"
	FunctionTypeVideosToVideo        = "videos_to_video"
)

// Workflow providers.
const (
	ProviderGPTImage1                       = "gpt-image-1"
	ProviderImagen3Generate002              = "imagen-3.0-generate-002"
	ProviderGemini20FlashExpImageGeneration = "gemini-2.0-flash-exp-image-generation"
	ProviderLeonardoKinoXL                  = "leonardo-kino-xl"
	ProviderLeonardoDiffusionXL             = "leonardo-diffusion-xl"
	ProviderLeonardoAnimeXL                 = "leonardo-anime-xl"
	ProviderLeonardoLightning               = "leonardo-lightning"
	ProviderDallE3                          = "dall-e-3"
	ProviderLumaPhoton                      = "luma/photon"
	ProviderLumaPhotonFlash                 = "luma/photon-flash"
	ProviderStabilitySD3                    = "stability-sd3"
	ProviderFluxSchnell                     = "flux-schnell"
	ProviderSana                            = "sana"
	ProviderSeedance1                       = replicate.Seedance1Model
	ProviderVeo3Preview                     = gemini.VEO_3_PREVIEW_MODEL
)

// WorkflowStep represents a single step in a workflow.
type WorkflowStep struct {
	ID           string   `json:"id" yaml:"id"`
	FunctionType string   `json:"function_type" yaml:"function_type"`
	Provider     string   `json:"provider,omitempty" yaml:"provider,omitempty"`
	Prompt       string   `json:"prompt,omitempty" yaml:"prompt,omitempty"`
	Image        string   `json:"image,omitempty" yaml:"image,omitempty"`
	FirstImage   string   `json:"first_image,omitempty" yaml:"first_image,omitempty"`
	LastImage    string   `json:"last_image,omitempty" yaml:"last_image,omitempty"`
	Videos       []string `json:"videos,omitempty" yaml:"videos,omitempty"`
}

// Workflow defines an ordered set of steps for content generation.
type Workflow struct {
	Name      string         `json:"name,omitempty" yaml:"name,omitempty"`
	Steps     []WorkflowStep `json:"steps" yaml:"steps"`
	CreatedAt time.Time      `json:"created_at" yaml:"created_at"`
	Output    string         `json:"output,omitempty" yaml:"output,omitempty"`
}

// WorkflowService executes workflows.
type WorkflowService interface {
	Generate(ctx context.Context, wf *Workflow, inputs map[string]any) (result any, output string, err error)
}

type workflowService struct{}

// NewWorkflowService returns a WorkflowService implementation.
func NewWorkflowService() WorkflowService {
	return &workflowService{}
}

// Generate executes a workflow with the provided inputs.
func (s *workflowService) Generate(ctx context.Context, wf *Workflow, inputs map[string]any) (any, string, error) {
	if wf == nil {
		return nil, "", errors.New("nil workflow")
	}

	results := make(map[string]any)
	for _, step := range wf.Steps {
		var (
			res any
			err error
		)
		switch step.FunctionType {
		case FunctionTypeTextsToText:
			res, err = s.processTextsToText(ctx, step, inputs, results)
		case FunctionTypeTextToImage:
			res, err = s.processTextToImage(ctx, step, inputs, results)
		case FunctionTypeTextAndImageToImage:
			res, err = s.processTextAndImageToImage(ctx, step, inputs, results)
		case FunctionTypeTextAndImagesToVideo:
			res, err = s.processTextAndImagesToVideo(ctx, step, inputs, results)
		case FunctionTypeTextAndImageToVideo:
			res, err = s.processTextAndImageToVideo(ctx, step, inputs, results)
		case FunctionTypeVideosToVideo:
			res, err = s.processVideosToVideo(ctx, step, inputs, results)
		default:
			err = errors.Errorf("unsupported function type: %s", step.FunctionType)
		}
		if err != nil {
			return nil, "", errors.Wrapf(err, "processing workflow step %s", step.ID)
		}
		results[step.ID] = res
	}

	var lastStepID string
	if len(wf.Steps) > 0 {
		lastStepID = wf.Steps[len(wf.Steps)-1].ID
	}
	var final any
	if lastStepID != "" {
		final = results[lastStepID]
	}
	return final, wf.Output, nil
}

func (s *workflowService) processTextsToText(ctx context.Context, step WorkflowStep, inputs map[string]any, results map[string]any) (any, error) {
	if step.Prompt == "" {
		return nil, errors.New("missing prompt template in step configuration")
	}

	out := s.interpolateVariables(step.Prompt, inputs, results)
	return out, nil
}

func (s *workflowService) processTextToImage(ctx context.Context, step WorkflowStep, inputs map[string]any, results map[string]any) (any, error) {
	// TODO: implement real logic. For now return dummy value.
	return "text_to_image result", nil
}

func (s *workflowService) processTextAndImageToImage(ctx context.Context, step WorkflowStep, inputs map[string]any, results map[string]any) (any, error) {
	// TODO: implement real logic. For now return dummy value.
	return "text_and_image_to_image result", nil
}

func (s *workflowService) processTextAndImagesToVideo(ctx context.Context, step WorkflowStep, inputs map[string]any, results map[string]any) (any, error) {
	if step.Prompt == "" {
		return nil, errors.New("missing prompt template in step configuration")
	}
	if step.FirstImage == "" || step.LastImage == "" {
		return nil, errors.New("missing first or last image in step configuration")
	}

	prompt := s.interpolateVariables(step.Prompt, inputs, results)
	firstURL := s.interpolateVariables(step.FirstImage, inputs, results)
	lastURL := s.interpolateVariables(step.LastImage, inputs, results)

	return s.generateVideo(ctx, step.Provider, prompt, firstURL, lastURL)
}

func (s *workflowService) processTextAndImageToVideo(ctx context.Context, step WorkflowStep, inputs map[string]any, results map[string]any) (any, error) {
	if step.Prompt == "" {
		return nil, errors.New("missing prompt template in step configuration")
	}
	if step.FirstImage == "" {
		return nil, errors.New("missing first image in step configuration")
	}

	prompt := s.interpolateVariables(step.Prompt, inputs, results)
	firstURL := s.interpolateVariables(step.FirstImage, inputs, results)

	return s.generateVideo(ctx, step.Provider, prompt, firstURL, "")
}

// generateVideo dispatches the video generation request to the chosen provider.
// If lastURL is empty, only the first frame is sent.
func (s *workflowService) generateVideo(ctx context.Context, provider, prompt, firstURL, lastURL string) (any, error) {
	if provider == "" {
		provider = ProviderVeo3Preview
	}

	switch provider {
	case ProviderVeo3Preview:
		svc := gemini.NewGeminiService()
		if lastURL != "" {
			return svc.GenerateVeo3PreviewVideoFromURLs(ctx, prompt, firstURL, lastURL)
		}
		return svc.GenerateVeo3PreviewVideoWithStartFrameURL(ctx, prompt, firstURL)
	case ProviderSeedance1:
		svc, err := replicate.NewReplicateService(os.Getenv(ReplicateAPIToken))
		if err != nil {
			return nil, err
		}
		opts := map[string]any{"image": firstURL}
		if lastURL != "" {
			opts["last_frame_image"] = lastURL
		}
		return svc.RunSeedance1(ctx, prompt, opts)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

func (s *workflowService) processVideosToVideo(ctx context.Context, step WorkflowStep, inputs map[string]any, results map[string]any) (any, error) {
	if len(step.Videos) == 0 {
		return nil, errors.New("no videos specified in step configuration")
	}

	var clips [][]byte
	for _, ref := range step.Videos {
		name := s.interpolateVariables(ref, inputs, results)
		var data any
		var ok bool
		if data, ok = results[name]; !ok {
			data, ok = inputs[name]
		}
		if !ok {
			return nil, fmt.Errorf("video reference %s not found", name)
		}
		b, ok := data.([]byte)
		if !ok {
			return nil, fmt.Errorf("video reference %s is not []byte", name)
		}
		clips = append(clips, b)
	}

	return MergeVideos(clips)
}

// interpolateVariables replaces placeholders in the template string with values from inputs and results.
func (s *workflowService) interpolateVariables(template string, inputs map[string]any, results map[string]any) string {
	result := template
	for k, v := range inputs {
		placeholder := fmt.Sprintf("${%s}", k)
		result = strings.ReplaceAll(result, placeholder, fmt.Sprint(v))
	}
	for k, v := range results {
		placeholder := fmt.Sprintf("${%s}", k)
		result = strings.ReplaceAll(result, placeholder, fmt.Sprint(v))
	}
	return result
}
