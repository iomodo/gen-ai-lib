package genailib

import (
	"context"
	"testing"
)

func TestWorkflowGenerate(t *testing.T) {
	svc := NewWorkflowService()
	wf := &Workflow{
		Steps: []WorkflowStep{
			{ID: "step1", FunctionType: FunctionTypeTextsToText},
			{ID: "step2", FunctionType: FunctionTypeTextToImage},
			{ID: "step3", FunctionType: FunctionTypeTextAndImageToImage},
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
