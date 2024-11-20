package boostrapper

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
	"google.golang.org/protobuf/proto"
)

type Bootstrapper struct {
	fileName string // json file with the neighbors of each node
	logger   *config.Logger
}

func NewBootstrapper(fileName string, logger *config.Logger) *Bootstrapper {
	return &Bootstrapper{
		fileName: fileName,
		logger:   logger,
	}
}

func (bs *Bootstrapper) Bootstrap(session quic.Connection, stream quic.Stream, header *protobuf.Header) {

  bs.logger.Info("Received Bootstrapper request.")

	// Read the JSON file
	data, err := os.ReadFile(bs.fileName)
	if err != nil {
		bs.logger.Error("Error reading file: " + err.Error())
		return
	}

	// Parse the JSON data
	var nodeNeighbors map[string][]string
	err = json.Unmarshal(data, &nodeNeighbors)
	if err != nil {
		bs.logger.Error("Error parsing JSON: " + err.Error())
		return
	}

  bs.logger.Info(fmt.Sprintf("Sending Neighbors: %v\n", nodeNeighbors[header.Sender]))

	message := protobuf.Header{
		Type:           protobuf.RequestType_RETRANSMIT,
		Length:         0,
		Timestamp:      time.Now().UnixMilli(),
		ClientIp:       "nil",
		Sender:         "bootstrapper",
		Target:         "nil",
		RequestedVideo: "nil",
		Content: &protobuf.Header_BootstraperResult{
			BootstraperResult: &protobuf.BootstraperResult{
				Neighbors: nodeNeighbors[header.Sender],
			},
		},
	}
  message.Length = int32(proto.Size(&message))
  message.Target = session.RemoteAddr().String()

  data, err = proto.Marshal(&message)

  if err != nil {
    bs.logger.Error(err.Error())
    return
  }

  config.SendMessage(stream, data)

	bs.logger.Info("Successfully sent node neighbors data to client.")
}
