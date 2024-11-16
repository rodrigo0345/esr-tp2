package bootstrapper

import (
	"fmt"
	"time"

	"github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
	"google.golang.org/protobuf/proto"
)

func BSGetNeighbors(cnf *config.AppConfigList, bsAddr *protobuf.Interface) ([]string, error) {
  message := protobuf.Header{
    Type:           protobuf.RequestType_BOOTSTRAPER,
    Length:         0,
    Timestamp:      time.Now().UnixMilli(),
    ClientIp:       "nil",
    Sender:         cnf.NodeName,
    Target:         "nil",
    RequestedVideo: "nil",
    Content: nil,
  }

  data, err := proto.Marshal(&message)
  if err != nil {
    return nil, err
  }

  addr := fmt.Sprintf("%s:%d", bsAddr.Ip, bsAddr.Port)

  // create a new connection
  stream, conn, err := config.StartConnStream(addr)
  if err != nil {
    return nil, err
  }
  defer config.CloseStream(stream)
  defer config.CloseConnection(conn)

  // send the message
  err = config.SendMessage(stream, data)
  if err != nil {
    return nil, err
  }

  // wait for the response
  dataRecv, err := config.ReceiveMessage(stream)
  if err != nil {
    return nil, err
  }
  proto.Unmarshal(dataRecv, &message)

  if message.GetType() != protobuf.RequestType_BOOTSTRAPER {
    return nil, err
  }

  neighbors := message.GetBootstraperResult().GetNeighbors()

  return neighbors, nil
}
