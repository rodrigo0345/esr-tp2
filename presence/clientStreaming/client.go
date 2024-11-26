package clientStreaming

import (
	"fmt"
	"net"

	"github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
	"google.golang.org/protobuf/proto"
)

func (ss *StreamingService) SendToClient(callback Callback, header *protobuf.Header) {

	mine := hasMessageForMe(header, ss.Config.NodeName)

	if !mine {
		callback <- CallbackData{
			Header: header,
			Cancel: true,
		}
		return
	}

	video := header.RequestedVideo
	otherNodesNeedThisPacket := len(header.GetTarget()) > 1

	for _, client := range ss.UdpClients[Video(video)] {
		data, err := proto.Marshal(header)
		if err != nil {
			ss.Logger.Error(err.Error())
			continue
		}

		udpAddr := fmt.Sprintf("%s:%d", client.Ip, client.Port)
    ss.Logger.Info(fmt.Sprintf("Sending via UDP to %s", udpAddr))
		go config.SendMessageUDP(udpAddr, data)
	}

	callback <- CallbackData{
		Header: header,
		Cancel: otherNodesNeedThisPacket,
	}
}

func hasMessageForMe(header *protobuf.Header, nodeName string) bool {
	target := header.GetTarget()

	for _, t := range target {
		if t == nodeName {
			return true
		}
	}
	return false
}

func SendViaUDP(data []byte, clientIp string) {
	const MaxPacketSize = 1400 // Choose a size less than the MTU
	udpAddr, err := net.ResolveUDPAddr("udp", clientIp)
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("Error dialing UDP connection:", err)
		return
	}
	defer conn.Close()

	for start := 0; start < len(data); start += MaxPacketSize {
		end := start + MaxPacketSize
		if end > len(data) {
			end = len(data)
		}
		_, err = conn.Write(data[start:end])
		if err != nil {
			fmt.Println("Error sending data:", err)
			return
		}
	}
}
