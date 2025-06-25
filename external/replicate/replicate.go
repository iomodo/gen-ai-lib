package replicate

import (
	"context"

	"github.com/pkg/errors"

	replicate "github.com/replicate/replicate-go"
)

// ReplicateServiceAPI defines the interface for Replicate service operations.
type ReplicateServiceAPI interface {
	Run(ctx context.Context, model string, input replicate.PredictionInput) (PredictionOutput, error)
}

// ReplicateService provides methods to interact with the Replicate API.
type ReplicateService struct {
	client *replicate.Client
}

// NewReplicateService returns a ReplicateServiceAPI interface.
func NewReplicateService(token string) (ReplicateServiceAPI, error) {
	client, err := replicate.NewClient(replicate.WithToken(token))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create replicate client")
	}
	return &ReplicateService{client: client}, nil
}

// Run uses the official Replicate Go client to execute a model prediction.
func (r *ReplicateService) Run(ctx context.Context, model string, input replicate.PredictionInput) (PredictionOutput, error) {
	output, err := r.client.Run(ctx, model, input, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to run model")
	}
	return output, nil
}

type PredictionOutput interface{}
