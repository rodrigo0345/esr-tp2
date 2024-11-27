package bootstrapsystem

import (
	"fmt"

	"github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
	"google.golang.org/protobuf/proto"
)

func (bs *BootstrapSystem) ReqPing(data *ChannelData) {

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
  }

  for _, nbAddr := range(bs.Neighbors) {
    bs.Logger.Info(fmt.Sprintf("Sending ping to neighbor %s", nbAddr))
    header.Target = []string{nbAddr}
    config.SendMessageUDP(nbAddr, d)
  }

	data.Callback <- &CallbackData{
		Success: true,
		Error:   nil,
	}
}
