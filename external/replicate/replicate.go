package replicate

import (
	"context"
	"maps"

	"github.com/pkg/errors"

	replicate "github.com/replicate/replicate-go"
)

// Seedance1Model identifies the bytedance/seedance-1 model on Replicate.
const Seedance1Model = "bytedance/seedance-1"

// Seedance1LiteModel identifies the bytedance/seedance-1-lite model on Replicate.
const Seedance1LiteModel = "bytedance/seedance-1-lite"

// ReplicateServiceAPI defines the interface for Replicate service operations.

type ReplicateService interface {
	Run(ctx context.Context, model string, prompt string, options map[string]any) (any, error)
	// RunSeedance1 runs the bytedance/seedance-1 model.
	// Options may include:
	//  - "image":               string or *replicate.File
	//  - "last_frame_image":    string or *replicate.File
	//  - "duration":           int (seconds)
	//  - "resolution":         string (e.g. "720p")
	//  - "aspect_ratio":       string (e.g. "16:9")
	//  - "fps":                int
	//  - "camera_fixed":       bool
	//  - "seed":               int
	RunSeedance1(ctx context.Context, prompt string, options map[string]any) (any, error)
	// RunSeedance1Lite runs the bytedance/seedance-1-lite model.
	// Options may include:
	//  - "image":               string or *replicate.File
	//  - "last_frame_image":    string or *replicate.File
	//  - "duration":           int (seconds)
	//  - "resolution":         string (e.g. "720p")
	//  - "aspect_ratio":       string (e.g. "16:9")
	//  - "fps":                int
	//  - "camera_fixed":       bool
	//  - "seed":               int
	RunSeedance1Lite(ctx context.Context, prompt string, options map[string]any) (any, error)
	IsInitialized() bool
}

// ReplicateService provides methods to interact with the Replicate API.
type replicateService struct {
	client *replicate.Client
}

// NewReplicateService returns a ReplicateServiceAPI interface.
func NewReplicateService(token string) (ReplicateService, error) {
	client, err := replicate.NewClient(replicate.WithToken(token))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create replicate client")
	}
	return &replicateService{client: client}, nil
}

// Run uses the official Replicate Go client to execute a model prediction.
func (r *replicateService) Run(ctx context.Context, model string, prompt string, options map[string]any) (any, error) {
	input := replicate.PredictionInput{
		"prompt": prompt,
	}
	maps.Copy(input, options)
	output, err := r.client.RunWithOptions(ctx, model, input, nil, replicate.WithBlockUntilDone())
	if err != nil {
		return nil, errors.Wrap(err, "failed to run model")
	}

	// Try to extract a URL from the output
	switch v := output.(type) {
	case []interface{}:
		if len(v) > 0 {
			if url, ok := v[0].(string); ok && (len(url) > 4 && (url[:4] == "http")) {
				return url, nil
			}
		}
	case map[string]interface{}:
		if url, ok := v["url"].(string); ok {
			return url, nil
		}
	}

	return output, nil
}

// RunSeedance1 executes the bytedance/seedance-1 model on Replicate.
// See the model's documentation for the supported input options.
func (r *replicateService) RunSeedance1(ctx context.Context, prompt string, options map[string]any) (any, error) {
	return r.Run(ctx, Seedance1Model, prompt, options)
}

// RunSeedance1Lite executes the bytedance/seedance-1-lite model on Replicate.
// See the model's documentation for the supported input options.
func (r *replicateService) RunSeedance1Lite(ctx context.Context, prompt string, options map[string]any) (any, error) {
	return r.Run(ctx, Seedance1LiteModel, prompt, options)
}

func (r *replicateService) IsInitialized() bool {
	return r.client != nil
}
