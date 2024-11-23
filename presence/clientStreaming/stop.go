package clientStreaming
import (
	"time"

	"github.com/rodrigo0345/esr-tp2/config/protobuf"
	"google.golang.org/protobuf/proto"
)

// this is directly called by the client
func (ss *StreamingService) Stop(video Video, udpClient *protobuf.Interface, callback Callback) {

  targets := make([]string, 1)
  targets[0] = "s1"

	msg := &protobuf.Header{
		Type:      protobuf.RequestType_RETRANSMIT,
		Length:    0,
    Target:    targets,
    Sender:    ss.Config.NodeName,
    Path:      ss.Config.NodeName,
		Timestamp: time.Now().UnixMilli(),
		Content: &protobuf.Header_ClientCommand{
			ClientCommand: &protobuf.ClientCommand{
				Command:               protobuf.PlayerCommand_STOP,
				AdditionalInformation: "no information",
			},
		},
		RequestedVideo: string(video),
	}
  msg.Length = int32(proto.Size(msg))

  // the video doesn't exist
	if _, exists := ss.UdpClients[video]; !exists {
    callback <- CallbackData {
      Header: msg,
      Cancel: true,
    }
    return 
	}

  client := UdpClient{Interface: udpClient, LastSeen: time.Now()}

  ss.RemoveUdpClient(video, client)

  callback <- CallbackData {
    Header: msg,
    Cancel: false,
  }
}
