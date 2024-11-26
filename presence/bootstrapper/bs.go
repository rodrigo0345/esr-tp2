package bootstrapper

import (
	"fmt"
	"time"

	"github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
	"google.golang.org/protobuf/proto"
)

func BSGetNeighbors(cnf *config.AppConfigList, bsAddr *protobuf.Interface) ([]string, error) {
  target := make([]string, 0)
  message := protobuf.Header{
    Type:           protobuf.RequestType_BOOTSTRAPER,
    Length:         0,
    Timestamp:      time.Now().UnixMilli(),
    Sender:         cnf.NodeName,
    Target:         target,
    RequestedVideo: "bootstrapper",
    Content: nil,
  }

  data, err := proto.Marshal(&message)
  if err != nil {
    return nil, err
  }

  addr := fmt.Sprintf("%s:%d", bsAddr.Ip, bsAddr.Port)

  // create a new connection
  stream, _, err := config.StartConnStream(addr)
  if err != nil {
    return nil, err
  }

  // send the message
  err = config.SendMessage(stream, data)
  if err != nil {
    return nil, err
  }
  defer stream.Close()

  // wait for the response
  dataRecv, err := config.ReceiveMessage(stream)
  if err != nil {
    return nil, err
  }
  proto.Unmarshal(dataRecv, &message)

  neighbors := message.GetBootstraperResult().GetNeighbors()

  return neighbors, nil
}
