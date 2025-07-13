package genailib

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
)

// AppendVideos takes two video byte slices and appends the second video to the first.
// It returns the merged video as a byte slice using ffmpeg under the hood.
func AppendVideos(video1, video2 []byte) ([]byte, error) {
	tmpDir, err := os.MkdirTemp("", "mergevideo")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create temp dir")
	}
	defer os.RemoveAll(tmpDir)

	input1 := filepath.Join(tmpDir, "input1.mp4")
	input2 := filepath.Join(tmpDir, "input2.mp4")
	output := filepath.Join(tmpDir, "output.mp4")

	if err := os.WriteFile(input1, video1, 0o600); err != nil {
		return nil, errors.Wrap(err, "failed to write first video")
	}
	if err := os.WriteFile(input2, video2, 0o600); err != nil {
		return nil, errors.Wrap(err, "failed to write second video")
	}

	cmd := exec.Command("ffmpeg", "-i", input1, "-i", input2, "-filter_complex", "[0:v][1:v]concat=n=2:v=1[out]", "-map", "[out]", "-y", output)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ffmpeg run error: %w, %s", err, stderr.String())
	}

	merged, err := os.ReadFile(output)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read merged video")
	}

	return merged, nil
}

// MergeVideos concatenates multiple MP4 video clips into a single video using
// ffmpeg. The input videos must all be encoded with compatible codecs for the
// concat demuxer to work correctly.
func MergeVideos(videos [][]byte) ([]byte, error) {
	if len(videos) == 0 {
		return nil, errors.New("no videos provided")
	}

	tmpDir, err := os.MkdirTemp("", "mergevideos")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create temp dir")
	}
	defer os.RemoveAll(tmpDir)

	listFile := filepath.Join(tmpDir, "inputs.txt")
	output := filepath.Join(tmpDir, "output.mp4")

	var list bytes.Buffer
	for i, data := range videos {
		filePath := filepath.Join(tmpDir, fmt.Sprintf("input%d.mp4", i))
		if err := os.WriteFile(filePath, data, 0o600); err != nil {
			return nil, errors.Wrapf(err, "failed to write video %d", i)
		}
		fmt.Fprintf(&list, "file '%s'\n", filePath)
	}

	if err := os.WriteFile(listFile, list.Bytes(), 0o600); err != nil {
		return nil, errors.Wrap(err, "failed to write list file")
	}

	cmd := exec.Command("ffmpeg", "-f", "concat", "-safe", "0", "-i", listFile, "-c", "copy", "-y", output)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ffmpeg run error: %w, %s", err, stderr.String())
	}

	merged, err := os.ReadFile(output)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read merged video")
	}

	return merged, nil
}
