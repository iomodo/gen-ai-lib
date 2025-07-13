package genailib

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func createColorVideo(color string) ([]byte, error) {
	tmpFile, err := os.CreateTemp("", color+"-*.mp4")
	if err != nil {
		return nil, err
	}
	tmpFile.Close()
	cmd := exec.Command("ffmpeg", "-f", "lavfi", "-i", fmt.Sprintf("color=c=%s:s=320x240:d=1", color), "-pix_fmt", "yuv420p", "-y", tmpFile.Name())
	if output, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("ffmpeg create video: %v, %s", err, string(output))
	}
	data, err := os.ReadFile(tmpFile.Name())
	os.Remove(tmpFile.Name())
	return data, err
}

func TestAppendVideos(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not installed")
	}
	v1, err := createColorVideo("red")
	if err != nil {
		t.Fatalf("failed to create first video: %v", err)
	}
	v2, err := createColorVideo("blue")
	if err != nil {
		t.Fatalf("failed to create second video: %v", err)
	}
	merged, err := AppendVideos(v1, v2)
	if err != nil {
		t.Fatalf("AppendVideos returned error: %v", err)
	}
	if len(merged) == 0 {
		t.Fatalf("merged video is empty")
	}
}

func TestMergeVideos(t *testing.T) {
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
	v3, err := createColorVideo("blue")
	if err != nil {
		t.Fatalf("failed to create third video: %v", err)
	}

	merged, err := MergeVideos([][]byte{v1, v2, v3})
	if err != nil {
		t.Fatalf("MergeVideos returned error: %v", err)
	}
	if len(merged) == 0 {
		t.Fatalf("merged video is empty")
	}
}
