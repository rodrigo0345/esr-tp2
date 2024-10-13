package server

import (
	"encoding/binary"
	"fmt"
	"time"
)

const (
	HEADER_SIZE = 12 // Size of the RTP header
)

// RTPPacket struct to hold the RTP header and payload
type RTPPacket struct {
	Header  []byte
	Payload []byte
}

// MakeRtpPacket constructs an RTP packet from the given payload and parameters.
func MakeRtpPacket(payload []byte, frameNumber uint16, start bool) *RTPPacket {
	version := 2
	padding := 0
	extension := 0
	cc := 0
	marker := 0
	if start {
		marker = 1 // Set marker if this is the start of a new frame
	}
	pt := 26 // MJPEG payload type
	ssrc := uint32(0)
	timestamp := uint32(time.Now().UnixNano() / 1000000) // Current timestamp in milliseconds

	rtpHeader := make([]byte, HEADER_SIZE)

	// Fill the RTP header fields
	rtpHeader[0] = (byte(version)<<6 | byte(padding)<<5 | byte(extension)<<4 | byte(cc&0x0F))
	rtpHeader[1] = (byte(marker)<<7 | byte(pt&0x7F))

	// Sequence Number
	rtpHeader[2] = byte(frameNumber >> 8)   // Upper 8 bits
	rtpHeader[3] = byte(frameNumber & 0xFF) // Lower 8 bits

	// Timestamp
	rtpHeader[4] = byte(timestamp >> 24)        // 1st byte
	rtpHeader[5] = byte(timestamp >> 16 & 0xFF) // 2nd byte
	rtpHeader[6] = byte(timestamp >> 8 & 0xFF)  // 3rd byte
	rtpHeader[7] = byte(timestamp & 0xFF)       // 4th byte

	// SSRC
	rtpHeader[8] = byte(ssrc >> 24)        // 1st byte of SSRC
	rtpHeader[9] = byte(ssrc >> 16 & 0xFF) // 2nd byte of SSRC
	rtpHeader[10] = byte(ssrc >> 8 & 0xFF) // 3rd byte of SSRC
	rtpHeader[11] = byte(ssrc & 0xFF)      // 4th byte of SSRC

	// Create and return the RTP packet
	return &RTPPacket{
		Header:  rtpHeader,
		Payload: payload,
	}
}

// Decode decodes the RTP packet from the byte stream.
func (r *RTPPacket) Decode(byteStream []byte) error {
	if len(byteStream) < HEADER_SIZE {
		return fmt.Errorf("byte stream is too short")
	}
	r.Header = make([]byte, HEADER_SIZE)
	copy(r.Header, byteStream[:HEADER_SIZE])
	r.Payload = byteStream[HEADER_SIZE:]
	return nil
}

// Version returns the RTP version.
func (r *RTPPacket) Version() int {
	return int(r.Header[0] >> 6)
}

// SeqNum returns the sequence number.
func (r *RTPPacket) SeqNum() uint16 {
	return binary.BigEndian.Uint16(r.Header[2:4])
}

// Timestamp returns the timestamp.
func (r *RTPPacket) Timestamp() uint32 {
	return binary.BigEndian.Uint32(r.Header[4:8])
}

// PayloadType returns the payload type.
func (r *RTPPacket) PayloadType() int {
	return int(r.Header[1] & 0x7F) // mask to get the lower 7 bits
}

// GetPayload returns the payload of the RTP packet.
func (r *RTPPacket) GetPayload() []byte {
	return r.Payload
}

// GetPacket returns the complete RTP packet (header + payload).
func (r *RTPPacket) GetPacket() []byte {
	return append(r.Header, r.Payload...)
}

// PrintHeader prints RTP header information.
func (r *RTPPacket) PrintHeader() {
	fmt.Printf("[RTP Packet] Version: %d, SeqNum: %d, Timestamp: %d, PayloadType: %d\n", r.Version(), r.SeqNum(), r.Timestamp(), r.PayloadType())
}
