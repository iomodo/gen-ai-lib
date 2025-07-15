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
	os.Setenv("GOOGLE_GENAI_USE_VERTEXAI", "true")
	if os.Getenv("GOOGLE_CLOUD_PROJECT") == "" {
		os.Setenv("GOOGLE_CLOUD_PROJECT", "test-project")
	}
	if os.Getenv("GOOGLE_CLOUD_LOCATION") == "" && os.Getenv("GOOGLE_CLOUD_REGION") == "" {
		os.Setenv("GOOGLE_CLOUD_LOCATION", "us-central1")
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

func TestWorkflowGenerateAndMergeVideos(t *testing.T) {
	if os.Getenv("GEMINI_API_KEY") == "" {
		t.Skip("GEMINI_API_KEY not set")
	}
	os.Setenv("GOOGLE_GENAI_USE_VERTEXAI", "true")
	if os.Getenv("GOOGLE_CLOUD_PROJECT") == "" {
		os.Setenv("GOOGLE_CLOUD_PROJECT", "test-project")
	}
	if os.Getenv("GOOGLE_CLOUD_LOCATION") == "" && os.Getenv("GOOGLE_CLOUD_REGION") == "" {
		os.Setenv("GOOGLE_CLOUD_LOCATION", "us-central1")
	}
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not installed")
	}

	svc := NewWorkflowService()
	wf := &Workflow{
		Steps: []WorkflowStep{
			{
				ID:           "video1",
				FunctionType: FunctionTypeTextAndImagesToVideo,
				Prompt:       "A quick clip one",
				FirstImage:   "https://picsum.photos/seed/frame1/256",
				LastImage:    "https://picsum.photos/seed/frame2/256",
			},
			{
				ID:           "video2",
				FunctionType: FunctionTypeTextAndImagesToVideo,
				Prompt:       "Another short clip",
				FirstImage:   "https://picsum.photos/seed/frame3/256",
				LastImage:    "https://picsum.photos/seed/frame4/256",
			},
			{
				ID:           "merge",
				FunctionType: FunctionTypeVideosToVideo,
				Videos:       []string{"video1", "video2"},
			},
		},
	}

	result, _, err := svc.Generate(context.Background(), wf, nil)
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

func TestWorkflowSingleImageClipsAndMerge(t *testing.T) {
	if os.Getenv("GEMINI_API_KEY") == "" {
		t.Skip("GEMINI_API_KEY not set")
	}
	os.Setenv("GOOGLE_GENAI_USE_VERTEXAI", "true")
	if os.Getenv("GOOGLE_CLOUD_PROJECT") == "" {
		os.Setenv("GOOGLE_CLOUD_PROJECT", "test-project")
	}
	if os.Getenv("GOOGLE_CLOUD_LOCATION") == "" && os.Getenv("GOOGLE_CLOUD_REGION") == "" {
		os.Setenv("GOOGLE_CLOUD_LOCATION", "us-central1")
	}
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not installed")
	}

	svc := NewWorkflowService()
	wf := &Workflow{
		Steps: []WorkflowStep{
			{
				ID:           "clip1",
				FunctionType: FunctionTypeTextAndImageToVideo,
				Prompt:       "Clip using one image",
				FirstImage:   "https://picsum.photos/seed/start1/256",
			},
			{
				ID:           "clip2",
				FunctionType: FunctionTypeTextAndImageToVideo,
				Prompt:       "Another single image clip",
				FirstImage:   "https://picsum.photos/seed/start2/256",
			},
			{
				ID:           "merge",
				FunctionType: FunctionTypeVideosToVideo,
				Videos:       []string{"clip1", "clip2"},
			},
		},
	}

	result, _, err := svc.Generate(context.Background(), wf, nil)
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

func TestWorkflowSeedanceSingleImageClipsAndMerge(t *testing.T) {
	if os.Getenv(ReplicateAPIToken) == "" {
		t.Skip("REPLICATE_API_TOKEN not set")
	}
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not installed")
	}

	svc := NewWorkflowService()
	wf := &Workflow{
		Steps: []WorkflowStep{
			{
				ID:           "clip1",
				FunctionType: FunctionTypeTextAndImageToVideo,
				Provider:     ProviderSeedance1Lite,
				Prompt:       "A seedance clip from one image",
				FirstImage:   "https://picsum.photos/seed/sclip1/256",
			},
			{
				ID:           "clip2",
				FunctionType: FunctionTypeTextAndImageToVideo,
				Provider:     ProviderSeedance1Lite,
				Prompt:       "Another seedance clip",
				FirstImage:   "https://picsum.photos/seed/sclip2/256",
			},
			{
				ID:           "merge",
				FunctionType: FunctionTypeVideosToVideo,
				Videos:       []string{"clip1", "clip2"},
			},
		},
	}

	result, _, err := svc.Generate(context.Background(), wf, nil)
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
