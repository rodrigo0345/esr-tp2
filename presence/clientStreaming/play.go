package clientStreaming

import (
	"time"

	"github.com/rodrigo0345/esr-tp2/config/protobuf"
	"google.golang.org/protobuf/proto"
)

// this is directly called by the client
func (ss *StreamingService) Play(video Video, udpClient *protobuf.Interface, callback Callback) {

	// Create a new stream for this video if it doesn't exist
	if _, exists := ss.UdpClients[video]; !exists {
		ss.UdpClients[video] = []UdpClient{}
	}

	// Verificar que o cliente ainda não está a receber o vídeo
	ss.Lock()
	for _, c := range ss.UdpClients[video] {
		if c.String() == udpClient.String() {
			ss.Unlock()
      ss.Logger.Error(ss.Config.NodeName + " is already streaming " + string(video) + " to " + udpClient.Ip)
      callback <- CallbackData {
        Header: nil,
        Cancel: true,
      }
      return
		}
	}
	ss.Unlock()

  client := UdpClient{Interface: udpClient, LastSeen: time.Now()}

	// Add the client to the list of clients interested in this video
	ss.AddUdpClient(video, client)

  targets := make([]string, 1)
  targets[0] = "s1"

	msg := &protobuf.Header{
		Type:      protobuf.RequestType_RETRANSMIT,
		Length:    0,
    Sender:    ss.Config.NodeName,
    Path:      ss.Config.NodeName,
		Timestamp: time.Now().UnixMilli(),
    Target:    targets,
		Content: &protobuf.Header_ClientCommand{
			ClientCommand: &protobuf.ClientCommand{
				Command:               protobuf.PlayerCommand_PLAY,
				AdditionalInformation: "no information",
			},
		},
		RequestedVideo: string(video),
	}
  msg.Length = int32(proto.Size(msg))

  callback <- CallbackData {
    Header: msg,
    Cancel: false,
  }
}
