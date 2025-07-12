package genailib

import (
	"errors"
	"log"
)

// ServiceFunc defines the function signature for API services.
type ServiceFunc func(args ...any) (any, error)

// service represents a single API endpoint.
type service struct {
	name string
	fn   ServiceFunc
}

// Standard errors used by the gateway.
var (
	ErrNoAPIAvailable    = errors.New("no API available")
	ErrContentPolicy     = errors.New("content policy violation")
	ErrInvalidParameters = errors.New("invalid parameters")
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
)

// APIGateway provides methods to select between different services.
type APIGateway struct{}

// NewAPIGateway returns a new APIGateway instance.
func NewAPIGateway() *APIGateway { return &APIGateway{} }

// Placeholder implementation for all service functions.
func notImplemented(args ...any) (any, error) { return nil, errors.New("not implemented") }

var (
	textToImageServices = []*service{
		{name: "imagen-3.0-generate-002", fn: notImplemented},
		{name: "gemini-2.0-flash-exp-image-generation", fn: notImplemented},
	}
)

func (g *APIGateway) getService(list []*service, preferredService string) *service {
	if preferredService != "" {
		for _, api := range list {
			if api.name == preferredService {
				return api
			}
		}
	}
	return textToImageServices[0]
}

// TextToImage executes a text-to-image workflow step.
func (g *APIGateway) TextToImage(prompt, preferredService string) (any, error) {
	api := g.getService(textToImageServices, preferredService)
	if api == nil {
		log.Printf("TextToImage error: %v", ErrNoAPIAvailable)
		return nil, ErrNoAPIAvailable
	}
	res, err := api.fn(prompt)
	if err != nil {
		if errors.Is(err, ErrContentPolicy) || errors.Is(err, ErrInvalidParameters) {
			return nil, err
		}
		if errors.Is(err, ErrRateLimitExceeded) {
			log.Printf("TextToImage %s rate limit exceeded", api.name)
		} else {
			log.Printf("TextToImage %s error: %v", api.name, err)
		}
		return nil, err
	}
	return res, nil
}
