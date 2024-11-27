package clientStreaming

import (
	"fmt"

	"github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
	"google.golang.org/protobuf/proto"
)

func (ss *StreamingService) ClientPing(myIp string, udpClient *protobuf.Interface, callback Callback) {

  msg := &protobuf.Header{
    Type: protobuf.RequestType_CLIENT_PING,
    Sender: myIp,
    Target: []string{udpClient.Ip},
  }

  data, err := proto.Marshal(msg)

  if err != nil {
    ss.Logger.Error(err.Error())
    return
  }

  config.SendMessageUDP(fmt.Sprintf("%s:%d", udpClient.Ip, udpClient.Port), data)

  callback <- CallbackData {
    Header: nil,
    Cancel: false,
  }
}
