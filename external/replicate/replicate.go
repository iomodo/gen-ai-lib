package replicate

import (
	"context"
	"maps"

	"github.com/pkg/errors"

	replicate "github.com/replicate/replicate-go"
)

// ReplicateServiceAPI defines the interface for Replicate service operations.
type ReplicateService interface {
	Run(ctx context.Context, model string, prompt string, options map[string]any) (any, error)
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

func (r *replicateService) IsInitialized() bool {
	return r.client != nil
}
