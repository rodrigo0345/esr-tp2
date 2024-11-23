package config

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/quic-go/quic-go"
)

var timeout time.Duration = 3 * time.Second // Set proper timeout duration

// StartConnStream initializes a QUIC connection and opens a stream
func StartConnStream(address string) (quic.Stream, quic.Connection, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}

	// Use context with timeout for dialing
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	connection, err := quic.DialAddr(ctx, address, tlsConfig, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to dial: %w", err)
	}

	// Open a stream with timeout
	ctx, cancel = context.WithTimeout(context.Background(), timeout)
	defer cancel()

	stream, err := connection.OpenStreamSync(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open stream: %w", err)
	}

	return stream, connection, nil
}

// CreateStream opens a new stream in the provided connection
func CreateStream(conn quic.Connection) (quic.Stream, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	stream, err := conn.OpenStreamSync(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to open stream: %w", err)
	}
	return stream, nil
}

// CloseStream closes a given stream
func CloseStream(stream quic.Stream) {
	if stream == nil {
		return
	}
	if err := stream.Close(); err != nil {
		log.Printf("Failed to close stream: %v", err)
	}
}

// CloseConnection closes the given QUIC connection
func CloseConnection(conn quic.Connection) {
  if conn == nil {
    return
  }
	if err := conn.CloseWithError(0, "closing connection"); err != nil {
		log.Printf("Failed to close connection: %v", err)
	}
}

// IsConnectionOpen checks if a connection is still open by writing a small message
func IsConnectionOpen(stream quic.Stream) bool {
	if stream == nil {
		return false
	}
	_, err := stream.Write([]byte("ping"))
	if err != nil {
		return false
	}
	return true
}

// SendMessage sends a message to the stream with length-prefix encoding
func SendMessage(stream quic.Stream, message []byte) error {
	if stream == nil {
		return fmt.Errorf("sending data to a nil stream")
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Send the length of the message first
	lengthPrefix := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthPrefix, uint32(len(message)))

	// Use a channel to handle the stream writing process
	writeDone := make(chan error, 1)
	go func() {
		defer close(writeDone)

		_, err := stream.Write(lengthPrefix)
		if err != nil {
			writeDone <- err
			return
		}

		_, err = stream.Write(message)
		writeDone <- err
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("send message timed out: %w", ctx.Err())
	case err := <-writeDone:
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
	}

	return nil
}

// ReceiveMessage reads a message from the stream with length-prefix decoding
func ReceiveMessage(stream quic.Stream) ([]byte, error) {
	if stream == nil {
		return nil, fmt.Errorf("receive message from nil stream")
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	lengthBuf := make([]byte, 4)

	// Read the length of the message
	readDone := make(chan error, 1)
	go func() {
		defer close(readDone)
		_, err := io.ReadFull(stream, lengthBuf)
		readDone <- err
	}()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("read message length timed out: %w", ctx.Err())
	case err := <-readDone:
		if err != nil {
			return nil, fmt.Errorf("failed to read message length: %w", err)
		}
	}

	messageLength := binary.BigEndian.Uint32(lengthBuf)
	if messageLength == 0 {
		return nil, fmt.Errorf("invalid message length: %d", messageLength)
	}

	buf := make([]byte, messageLength)

	// Read the actual message
	readDone = make(chan error, 1)
	go func() {
		defer close(readDone)
		_, err := io.ReadFull(stream, buf)
		readDone <- err
	}()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("read message timed out: %w", ctx.Err())
	case err := <-readDone:
		if err != nil {
			return nil, fmt.Errorf("failed to read message: %w", err)
		}
	}

	return buf, nil
}

