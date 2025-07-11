package genailib

import (
	"context"
	"os"

	"github.com/iomodo/gen-ai-lib/external/replicate"
	"github.com/pkg/errors"
)

const (
	ReplicateProvider = "replicate"
)

const (
	ReplicateAPIToken = "REPLICATE_API_TOKEN"
)

// Image defines the interface for generative image models.
type Image interface {
	// Generate creates a new image based on the given prompt and options.
	Generate(provider string, model string, prompt string, options map[string]interface{}) (any, error)

	// Edit modifies an existing image based on the given prompt and options.
	Edit(provider string, model string, input any, prompt string, options map[string]interface{}) (any, error)
}

type image struct {
	replicateService replicate.ReplicateService
}

func NewImageService() Image {
	return &image{}
}

func (i *image) Generate(provider string, model string, prompt string, options map[string]any) (any, error) {
	if provider == ReplicateProvider {
		replicateService, err := i.getReplicateService()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get replicate service")
		}
		output, err := replicateService.Run(context.Background(), model, prompt, options)
		if err != nil {
			return nil, errors.Wrap(err, "failed to run replicate model")
		}
		return output, nil
	}
	return nil, nil
}

func (i *image) Edit(provider string, model string, input any, prompt string, options map[string]any) (any, error) {
	return nil, nil
}

func (i *image) getReplicateService() (replicate.ReplicateService, error) {
	if i.replicateService == nil || !i.replicateService.IsInitialized() {
		replicateService, err := replicate.NewReplicateService(os.Getenv(ReplicateAPIToken))
		if err != nil {
			return nil, err
		}
		i.replicateService = replicateService
		return i.replicateService, nil
	}
	return i.replicateService, nil
}
