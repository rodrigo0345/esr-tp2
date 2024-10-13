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

func (vs *VideoStream) NextFrame() ([]byte, error) {
	// Read the first 5 bytes to get the frame length (as a string)
	frameLengthBytes := make([]byte, 5)
	_, err := io.ReadFull(vs.file, frameLengthBytes)
	if err != nil {
		if err == io.EOF {
			// Restart the stream when EOF is reached and loop through again
			vs.Restart()
			return vs.NextFrame()
		}
		return nil, fmt.Errorf("error reading frame length: %w", err)
	}

	// Convert the frame length string to an integer
	frameLength, err := strconv.Atoi(string(frameLengthBytes))
	if err != nil {
		return nil, fmt.Errorf("invalid frame length: %w", err)
	}

	// Add a validation step to avoid corrupt frames
	if frameLength <= 0 || frameLength > 500000 { // Adjust max size as per expected frame sizes
		return nil, fmt.Errorf("frame length out of bounds: %d", frameLength)
	}

	// Read the frame data based on the frame length
	frameData := make([]byte, frameLength)
	_, err = io.ReadFull(vs.file, frameData)
	if err != nil {
		if err == io.EOF {
			// If EOF, attempt to restart the stream and retrieve next frame
      // This loops the stream
			vs.Restart()
			return vs.NextFrame()
		}
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
