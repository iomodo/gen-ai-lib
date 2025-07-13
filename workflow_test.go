package genailib

import (
	"context"
	"os"
	"testing"
)

func TestWorkflowGenerate(t *testing.T) {
	svc := NewWorkflowService()
	wf := &Workflow{
		Steps: []WorkflowStep{
			{ID: "step1", FunctionType: FunctionTypeTextsToText, Prompt: "hello"},
			{ID: "step2", FunctionType: FunctionTypeTextToImage, Prompt: "image"},
			{ID: "step3", FunctionType: FunctionTypeTextAndImageToImage, Prompt: "edit"},
		},
	}
	result, _, err := svc.Generate(context.Background(), wf, nil)
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	if result != "text_and_image_to_image result" {
		t.Fatalf("unexpected result: %v", result)
	}
}

func TestWorkflowGenerateVideo(t *testing.T) {
	if os.Getenv("GEMINI_API_KEY") == "" {
		t.Skip("GEMINI_API_KEY not set")
	}
	svc := NewWorkflowService()
	wf := &Workflow{
		Steps: []WorkflowStep{
			{
				ID:           "step1",
				FunctionType: FunctionTypeTextAndImagesToVideo,
				Prompt:       "Create a short clip of a cat",
				FirstImage:   "https://picsum.photos/seed/cat1/256",
				LastImage:    "https://picsum.photos/seed/cat2/256",
			},
		},
	}
	_, _, err := svc.Generate(context.Background(), wf, nil)
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
}

func TestWorkflowGenerateUnsupported(t *testing.T) {
	svc := NewWorkflowService()
	wf := &Workflow{Steps: []WorkflowStep{{ID: "step1", FunctionType: "unknown"}}}
	_, _, err := svc.Generate(context.Background(), wf, nil)
	if err == nil {
		t.Fatal("expected error for unsupported function type")
	}
}

func TestWorkflowGenerateNil(t *testing.T) {
	svc := NewWorkflowService()
	_, _, err := svc.Generate(context.Background(), nil, nil)
	if err == nil {
		t.Fatal("expected error for nil workflow")
	}
}
