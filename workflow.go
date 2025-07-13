package genailib

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// WorkflowStepFunctionType constants.
const (
	FunctionTypeTextsToText         = "texts_to_text"
	FunctionTypeTextToImage         = "text_to_image"
	FunctionTypeTextAndImageToImage = "text_and_image_to_image"
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
)

// WorkflowStep represents a single step in a workflow.
type WorkflowStep struct {
	ID           string `json:"id" yaml:"id"`
	FunctionType string `json:"function_type" yaml:"function_type"`
	Provider     string `json:"provider,omitempty" yaml:"provider,omitempty"`
	Prompt       string `json:"prompt,omitempty" yaml:"prompt,omitempty"`
	Image        string `json:"image,omitempty" yaml:"image,omitempty"`
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
