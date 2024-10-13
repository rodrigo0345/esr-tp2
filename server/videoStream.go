package server

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

type VideoStream struct {
	filename string
	file     *os.File
	frameNum int
}

// NewVideoStream creates a new VideoStream instance
func NewVideoStream(filename string) (*VideoStream, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return &VideoStream{
		filename: filename,
		file:     file,
		frameNum: 0,
	}, nil
}

// NextFrame retrieves the next frame from the video stream
func (vs *VideoStream) NextFrame() ([]byte, error) {
	// Read the first 5 bytes to get the frame length (as a string)
	frameLengthBytes := make([]byte, 5)
	_, err := io.ReadFull(vs.file, frameLengthBytes)
	if err != nil {
		if err == io.EOF {
			vs.Restart()
			return vs.NextFrame() // Restart and try again
		}
		return nil, fmt.Errorf("error reading frame length: %w", err)
	}

	// Convert the frame length string to an integer
	frameLength, err := strconv.Atoi(string(frameLengthBytes))
	if err != nil {
		return nil, fmt.Errorf("invalid frame length: %w", err)
	}

	// Read the frame data based on the frame length
	frameData := make([]byte, frameLength)
	_, err = io.ReadFull(vs.file, frameData)
	if err != nil {
		return nil, fmt.Errorf("error reading frame data: %w", err)
	}

	vs.frameNum++
	return frameData, nil
}

// Restart resets the stream to the beginning of the file
func (vs *VideoStream) Restart() error {
	_, err := vs.file.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("failed to restart stream: %w", err)
	}
	// vs.frameNum = 0 // Uncomment if frame number needs to reset
	return nil
}

// FrameNumber returns the current frame number
func (vs *VideoStream) FrameNumber() int {
	return vs.frameNum
}

// Close closes the file stream
func (vs *VideoStream) Close() error {
	return vs.file.Close()
}
