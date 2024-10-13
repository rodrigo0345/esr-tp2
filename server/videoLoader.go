package server

import (
	"fmt"
	"os"
)

// it is up to the caller to close the returned file
func loadVideo(videoUrl string) (*os.File, error) {
	file, err := os.Open(videoUrl)
	if err != nil {
		return nil, fmt.Errorf("Failed to open video file: %v", err)
	}
	return file, nil
}
