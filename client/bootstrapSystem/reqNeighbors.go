package bootstrapsystem

import (
	"github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
	"google.golang.org/protobuf/proto"
)

func (bs *BootstrapSystem) ReqNeighbors(data *ChannelData) {
  serverIp := data.ContactAddr

  header := &protobuf.Header{
    Type: protobuf.RequestType_CLIENT_PING,
    Sender: "client",
  }

  d, err := proto.Marshal(header)
  if err != nil {
    data.Callback <- &CallbackData{
      Success: false,
      Error:   err,
    }
    return
  }

  err = config.SendMessageUDP(serverIp, d)
  if err != nil {
    data.Callback <- &CallbackData{
      Success: false,
      Error:   err,
    }
    return
  }

  bs.Logger.Info("Sending request to bootstraper")
  data.Callback <- &CallbackData{
    Success: true,
    Error:   nil,
  }
}
