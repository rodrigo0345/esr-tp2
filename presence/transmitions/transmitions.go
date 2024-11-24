package transmitions

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
	dvr "github.com/rodrigo0345/esr-tp2/presence/distance_vector_routing"
	"google.golang.org/protobuf/proto"
)

type TransmissionService struct {
	Logger       *config.Logger
	Config       *config.AppConfigList
}

func NewTransmissionService(logger *config.Logger, config *config.AppConfigList) *TransmissionService {
	return &TransmissionService{
		Logger:       logger,
		Config:       config,
	}
}

func (ts *TransmissionService) SendPacket(packet *protobuf.Header, rt *dvr.DistanceVectorRouting, isUdp bool) bool {

  if isUdp {
    ts.Logger.Debug("Using UDP")
  }
	messages := SplitTargets(packet.Target, rt, packet)

	if len(messages) == 0 {
		ts.Logger.Error("No targets found")
		return false
	}

	for node, packet := range messages {

		// sometimes it happens that the node sends the packet to itself
		if node.Ip == ts.Config.NodeIP.Ip {
			continue
		}

		neighborIp := fmt.Sprintf("%s:%d", node.Ip, node.Port)
		data, err := proto.Marshal(packet)

		if err != nil {
			ts.Logger.Error(err.Error())
			continue
		}

    if !isUdp {
      go sendData(data, neighborIp)
    } else {
      go sendDataUDP(data, neighborIp)
    }
	}
	return true
}

func sendData(data []byte, clientIp string) {

	neighborStream, conn, err := config.StartConnStream(clientIp)
	if err != nil {
		goto fail
	}

	err = config.SendMessage(neighborStream, data)

	if err != nil {
		goto fail
	}
	return

fail:
	config.CloseStream(neighborStream)
	config.CloseConnection(conn)
}

func sendDataUDP(data []byte, clientIp string) {
  // split the ip and port
  port := strings.Split(clientIp, ":")[1]
  ip := strings.Split(clientIp, ":")[0]

  portInt, err := strconv.Atoi(port)
  if err != nil {
    fmt.Println(err)
    return
  }

  // update the port
  clientIp = fmt.Sprintf("%s:%d", ip, portInt - 1)
  fmt.Println(fmt.Sprintf("Sending UDP packet to %s", clientIp))

  config.SendMessageUDP(clientIp, data)
}

// this enables us to not have duplicate streams of the same video
func SplitTargets(targets []string, routingTable *dvr.DistanceVectorRouting, packet *protobuf.Header) map[*protobuf.Interface]*protobuf.Header {
	routingTable.Mutex.RLock()
	defer routingTable.Mutex.RUnlock()

	// Initialize a map to store targets grouped by next hop.
	nextHops := make(map[*protobuf.Interface][]string)

	// Loop through each target to determine its next hop.
	for _, target := range targets {
		nextHop, err := routingTable.GetNextHop(target)
		if err != nil || nextHop.NextNode == nil { // Handle cases where no valid next hop is found.
			continue
		}
		// Group targets by next hop.
		nextHops[nextHop.NextNode] = append(nextHops[nextHop.NextNode], target)
	}

	// Initialize the result map.
	result := make(map[*protobuf.Interface]*protobuf.Header)

	// Create a protobuf.Header for each next hop with its associated targets.
	for node, targets := range nextHops {
		// Ensure each header is independently allocated.
		header := &protobuf.Header{
			Type:           packet.Type,
			Length:         packet.Length,
			Timestamp:      packet.Timestamp,
			Sender:         packet.Sender,
			Target:         targets, // Assign the specific targets for this next hop.
			RequestedVideo: packet.RequestedVideo,
			Content:        packet.Content,
			Path:           packet.Path,
		}
		result[node] = header
	}

	return result
}
