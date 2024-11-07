package config

import (
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
	"google.golang.org/protobuf/proto"
)

func UnmarshalHeader(data []byte) (*protobuf.Header, error) {
	chunk := &protobuf.Header{}
	err := proto.Unmarshal(data, chunk)
	if err != nil {
		return nil, err
	}
	return chunk, nil
}

func MarshalHeader(chunk *protobuf.Header) ([]byte, error) {
	data, err := proto.Marshal(chunk)
	if err != nil {
		return nil, err
	}
	return data, nil
}
