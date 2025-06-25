package genailib

// Image defines the interface for generative image models.
type Image interface {
	// Generate creates a new image based on the given prompt and options.
	Generate(provider string, model string, prompt string, options ...Option) (ImageResult, error)

	// Edit modifies an existing image based on the given prompt and options.
	Edit(provider string, model string, input ImageInput, prompt string, options ...Option) (ImageResult, error)
}

// Option represents an optional parameter for image generation or editing.
type Option interface{}

// ImageInput represents the input for image editing (could be a file, URL, or bytes).
type ImageInput interface{}

// ImageResult represents the result of an image generation or edit operation.
type ImageResult interface{}
