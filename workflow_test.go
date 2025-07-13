package genailib

import (
	"context"
	"os"
	"os/exec"
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

func TestWorkflowMergeVideos(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not installed")
	}

	v1, err := createColorVideo("red")
	if err != nil {
		t.Fatalf("failed to create first video: %v", err)
	}
	v2, err := createColorVideo("green")
	if err != nil {
		t.Fatalf("failed to create second video: %v", err)
	}

	svc := NewWorkflowService()
	wf := &Workflow{
		Steps: []WorkflowStep{
			{ID: "step1", FunctionType: FunctionTypeVideosToVideo, Videos: []string{"clip1"}},
			{ID: "step2", FunctionType: FunctionTypeVideosToVideo, Videos: []string{"clip2"}},
			{
				ID:           "merge",
				FunctionType: FunctionTypeVideosToVideo,
				Videos:       []string{"step1", "step2"},
			},
		},
	}

	inputs := map[string]any{
		"clip1": v1,
		"clip2": v2,
	}

	result, _, err := svc.Generate(context.Background(), wf, inputs)
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	merged, ok := result.([]byte)
	if !ok {
		t.Fatalf("expected []byte result, got %T", result)
	}
	if len(merged) == 0 {
		t.Fatalf("merged video is empty")
	}
}
