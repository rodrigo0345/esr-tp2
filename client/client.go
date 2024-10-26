package client

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	"io"
	"log"
	"os"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf" // Import your generated protobuf package
	"google.golang.org/protobuf/proto"
)

func Client(config *config.AppConfigList) {
	// Create a context for the connection
	ctx := context.Background()

	// Dial the QUIC server with the proper signature
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // Skip verification for testing purposes
		NextProtos:         []string{"quic-echo-example"},
	}
	quicConfig := &quic.Config{} // You can configure QUIC settings here, or leave it nil

	// Dial the QUIC server
	session, err := quic.DialAddr(ctx, "localhost:4242", tlsConfig, quicConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Open a stream to send and receive data
	stream, err := session.OpenStreamSync(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Open the video file
	videoFile, err := os.Open("videos/lol.Mjpeg")
	if err != nil {
		log.Fatal(err)
	}
	defer videoFile.Close()

	// Create a buffer to hold chunks of the video
	buffer := make([]byte, 1024)            // Adjust chunk size as necessary
	sequenceNumber := 1                     // Initialize sequence number
	videoFormat := protobuf.VideoFormat_MP4 // Set the video format (adjust accordingly)

	for {
		n, err := videoFile.Read(buffer)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		if n == 0 {
			break
		}

		videoChunk := &protobuf.ServerVideoChunk{
			SequenceNumber: int32(sequenceNumber),
			Timestamp:      time.Now().UnixMilli(),
			Format:         videoFormat,
			Data:           buffer[:n],
			IsLastChunk:    false,
		}

		serializedChunk, err := proto.Marshal(videoChunk)
		if err != nil {
			log.Fatal("Failed to serialize chunk:", err)
		}

		// Send the length of the message first
		lengthPrefix := make([]byte, 4)
		binary.BigEndian.PutUint32(lengthPrefix, uint32(len(serializedChunk)))

		// Send the length prefix followed by the serialized chunk
		_, err = stream.Write(lengthPrefix)
		if err != nil {
			log.Fatal("Failed to send length prefix:", err)
		}
		_, err = stream.Write(serializedChunk)
		if err != nil {
			log.Fatal("Failed to send chunk:", err)
		}

		sequenceNumber++
	}
}
