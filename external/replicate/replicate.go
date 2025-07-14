package replicate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
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
	token    string
	client   *http.Client
	baseURL  string
	versions map[string]string
}

// NewReplicateService returns a ReplicateServiceAPI interface using the HTTP API.
func NewReplicateService(token string) (ReplicateService, error) {
	if token == "" {
		return nil, errors.New("replicate token is required")
	}
	return &replicateService{
		token:    token,
		client:   &http.Client{Timeout: 60 * time.Second},
		baseURL:  "https://api.replicate.com/v1",
		versions: map[string]string{},
	}, nil
}

// Run executes a model prediction using Replicate's HTTP API and waits for completion.
func (r *replicateService) Run(ctx context.Context, model string, prompt string, options map[string]any) (any, error) {
	version, err := r.getLatestVersion(ctx, model)
	if err != nil {
		return nil, err
	}

	input := map[string]any{"prompt": prompt}
	maps.Copy(input, options)

	body, err := json.Marshal(map[string]any{
		"version": version,
		"input":   input,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.baseURL+"/predictions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Token "+r.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create prediction")
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("prediction create failed: %s", string(b))
	}

	var pred prediction
	if err := json.NewDecoder(resp.Body).Decode(&pred); err != nil {
		return nil, err
	}

	// Poll until prediction finished
	for {
		if pred.Status == "succeeded" {
			return extractOutputURL(pred.Output), nil
		}
		if pred.Status == "failed" || pred.Status == "canceled" {
			return nil, fmt.Errorf("prediction %s", pred.Status)
		}

		if err := r.getPrediction(ctx, pred.ID, &pred); err != nil {
			return nil, err
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			time.Sleep(2 * time.Second)
		}
	}
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

type prediction struct {
	ID     string      `json:"id"`
	Status string      `json:"status"`
	Output interface{} `json:"output"`
}

func (r *replicateService) getPrediction(ctx context.Context, id string, pred *prediction) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.baseURL+"/predictions/"+id, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Token "+r.token)
	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("prediction fetch failed: %s", string(b))
	}
	return json.NewDecoder(resp.Body).Decode(pred)
}

func (r *replicateService) getLatestVersion(ctx context.Context, model string) (string, error) {
	if v, ok := r.versions[model]; ok {
		return v, nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/models/%s", r.baseURL, model), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Token "+r.token)
	resp, err := r.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("fetch model failed: %s", string(b))
	}

	var info struct {
		LatestVersion struct {
			ID string `json:"id"`
		} `json:"latest_version"`
		DefaultVersion struct {
			ID string `json:"id"`
		} `json:"default_version"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return "", err
	}
	version := info.DefaultVersion.ID
	if version == "" {
		version = info.LatestVersion.ID
	}
	if version == "" {
		return "", errors.New("model version not found")
	}
	r.versions[model] = version
	return version, nil
}

func extractOutputURL(output interface{}) interface{} {
	switch v := output.(type) {
	case []interface{}:
		if len(v) > 0 {
			if url, ok := v[0].(string); ok && strings.HasPrefix(url, "http") {
				return url
			}
		}
	case map[string]interface{}:
		if url, ok := v["url"].(string); ok {
			return url
		}
	}
	return output
}
