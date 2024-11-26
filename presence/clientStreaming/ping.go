package clientStreaming

import (
	"time"

	"github.com/rodrigo0345/esr-tp2/config/protobuf"
)

// this is directly called by the client
func (ss *StreamingService) ClientPing(video Video, udpClient *protobuf.Interface, callback Callback) {

  client := UdpClient{Interface: udpClient, LastSeen: time.Now()}
  ss.MarkClientAsSeen(client)

  callback <- CallbackData {
    Header: nil,
    Cancel: false,
  }
}
