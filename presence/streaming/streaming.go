package streaming

import "github.com/quic-go/quic-go"

type RequestStream struct {
  conn quic.Connection
	stream quic.Stream
  source string
  video string
}

type ResponseStream struct {
  conn quic.Connection
  stream quic.Stream
}


