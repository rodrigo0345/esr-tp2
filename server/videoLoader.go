package server

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// it is up to the caller to close the returned file
func loadVideo(videoUrl string) (*os.File, error) {
	file, err := os.Open(videoUrl)
	if err != nil {
		return nil, fmt.Errorf("Failed to open video file: %v", err)
	}
	return file, nil
}

func ConvertVideoToH264(input string, output string) error {
	// get input file extension
	ext := strings.Replace(filepath.Ext(input), ".", "", -1)

	// Prepare the command to scale the video if needed
  cmd := exec.Command("ffmpeg", "-y", "-i", input, "-vf", "scale=iw-mod(iw\\,2):ih-mod(ih\\,2)", "-b:v", "1M", "-f", ext, output)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
