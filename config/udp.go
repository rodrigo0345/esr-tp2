package config

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/sys/unix"
)

// SendMessageUDP sends a message to the given UDP address
func SendMessageUDP(address string, message []byte) error {
	conn, err := net.Dial("udp", address)
	if err != nil {
		return fmt.Errorf("failed to dial UDP address %s: %w", address, err)
	}
	defer conn.Close()

	// Send the message directly
	_, err = conn.Write(message)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

// SetNonBlocking sets the given UDP connection to non-blocking mode
func SetNonBlocking(conn *net.UDPConn) error {
	rawConn, err := conn.SyscallConn()
	if err != nil {
		return fmt.Errorf("failed to get raw connection: %w", err)
	}

	return rawConn.Control(func(fd uintptr) {
		unix.SetNonblock(int(fd), true)
	})
}

// SetUDPBufferSize sets the buffer size for a UDP connection
func SetUDPBufferSize(conn *net.UDPConn, size int) error {
	rawConn, err := conn.SyscallConn()
	if err != nil {
		return fmt.Errorf("failed to get raw connection: %w", err)
	}

	return rawConn.Control(func(fd uintptr) {
		unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_RCVBUF, size)
	})
}

// ReceiveMessageUDP listens on the given UDP connection and receives a single message
func ReceiveMessageUDP(conn *net.UDPConn) ([]byte, *net.UDPAddr, error) {
	// Allocate a buffer to hold the incoming datagram
	buf := make([]byte, 65535)

	// Set a read deadline to avoid indefinite blocking
	conn.SetReadDeadline(time.Now().Add(50 * time.Millisecond))

	// Non-blocking read
	n, addr, err := conn.ReadFromUDP(buf)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			// Timeout occurred; this indicates no data is available
			return nil, nil, nil
		}
		return nil, nil, fmt.Errorf("failed to read message: %w", err)
	}

	if n == 0 {
		return nil, nil, fmt.Errorf("received empty message")
	}

	// Return the slice containing the actual message data
	return buf[:n], addr, nil
}

