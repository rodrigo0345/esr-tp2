package clientStreaming

import (
	"fmt"
	"sync"
	"time"

	"github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
	"google.golang.org/protobuf/proto"
)

type Video string
type Node string
type Command string
type CallbackData struct {
	Header *protobuf.Header
	Cancel bool
}
type Callback chan CallbackData

type SignalData struct {
	Command   Command
	Video     Video
	UdpClient *protobuf.Interface
	Callback  Callback
	Packet    *protobuf.Header
	MyIp      string
}

type UdpClient struct {
	*protobuf.Interface
	LastSeen time.Time
}

type Signal chan SignalData

const (
	PLAY      Command = "play"
	STOP      Command = "stop"
	HEARTBEAT Command = "heartbeat"
	VIDEO     Command = "video"
	PING      Command = "ping"
)

type StreamingService struct {
	sync.Mutex
	UdpClients      map[Video][]*UdpClient
	InterestedNodes map[Video][]Node
	Logger          *config.Logger
	Config          *config.AppConfigList
	Signal          Signal
	SignalDead      Callback
}

func (ss *StreamingService) RunBackgroundRoutine(seconds time.Duration) {

	go func() {
		for {
			select {
			case signal := <-ss.Signal:
				switch signal.Command {
				case PLAY:
					ss.Play(signal.Video, signal.UdpClient, signal.Callback)
				case STOP:
					ss.Stop(signal.Video, signal.UdpClient, signal.Callback)
				case HEARTBEAT:
					ss.ClientHeartbeat(signal.Video, signal.UdpClient, signal.Callback)
				case VIDEO:
					ss.SendToClient(signal.Callback, signal.Packet)
				case PING:
					ss.ClientPing(signal.MyIp, signal.UdpClient, signal.Callback)
				}
			}
		}
	}()

	// check inactive clients
	go func() {
		for {
			ss.CheckInactiveClients(seconds)
			time.Sleep(time.Second * 15)
		}
	}()
}

func (ss *StreamingService) SendSignal(signal SignalData) {
	ss.Signal <- signal
}

func (ss *StreamingService) Lock() {
	ss.Mutex.Lock()
}

func (ss *StreamingService) Unlock() {
	ss.Mutex.Unlock()
}

func NewStreamingService(cnf *config.AppConfigList, logger *config.Logger) *StreamingService {
	return &StreamingService{
		UdpClients:      make(map[Video][]*UdpClient),
		InterestedNodes: make(map[Video][]Node),
		Logger:          logger,
		Config:          cnf,
		SignalDead:      make(Callback),
		Signal:          make(Signal),
	}
}

func (ss *StreamingService) AddUdpClient(video Video, client *UdpClient) {
	ss.UdpClients[video] = append(ss.UdpClients[video], client)
}

func (ss *StreamingService) CheckInactiveClients(seconds time.Duration) {

	for video, clients := range ss.UdpClients {
		for _, client := range clients {
			if time.Since(client.LastSeen) > seconds*time.Second {

				ss.Logger.Error(fmt.Sprintf("Last seen %s:%d in %s since %d seconds", client.Ip, client.Port, video, seconds))
				ss.RemoveUdpClient(video, *client)

        if ss.UdpClients[video] != nil {
          continue
        }

        // notifies the network
        ss.SignalDead <- CallbackData{
          Header: ss.PrepareMessageForDeadClient(video, *client),
          Cancel: false,
        }

			}
		}
	}
}

func (ss *StreamingService) PrepareMessageForDeadClient(video Video, client UdpClient) *protobuf.Header {
	target := make([]string, 1)
	target[0] = "s1"
	header := &protobuf.Header{
		Type:           protobuf.RequestType_RETRANSMIT,
		Length:         0,
		Timestamp:      time.Now().UnixMilli(),
		Sender:         ss.Config.NodeName,
		Target:         target, // TODO: change this and make it dynamic
		RequestedVideo: string(video),
		Content: &protobuf.Header_ClientCommand{
			ClientCommand: &protobuf.ClientCommand{
				Command:               protobuf.PlayerCommand_STOP,
				AdditionalInformation: "client no longer pinged",
			},
		},
	}
	header.Length = int32(proto.Size(header))
	return header
}

func (ss *StreamingService) MarkClientAsSeen(client UdpClient) {
	ss.Lock()
	defer ss.Unlock()

	for _, clients := range ss.UdpClients {
		for _, cl := range clients {

			clString := fmt.Sprintf("%s:%d", cl.Ip, cl.Port)
			clientString := fmt.Sprintf("%s:%d", client.Ip, client.Port)

			if clString == clientString {
				cl.LastSeen = time.Now()
			}
		}
	}
}

func (ss *StreamingService) RemoveUdpClient(video Video, client UdpClient) {
	ss.Lock()
	defer ss.Unlock()

	for i, c := range ss.UdpClients[video] {
		if c.String() == client.String() {
			ss.UdpClients[video] = append(ss.UdpClients[video][:i], ss.UdpClients[video][i+1:]...)
			break
		}
	}

	// if the list is empty, remove the video
	if len(ss.UdpClients[video]) == 0 {
		delete(ss.UdpClients, video)
	}
}

func (ss *StreamingService) AddInterestedNode(video Video, node Node) {
	ss.Lock()
	defer ss.Unlock()
	ss.InterestedNodes[video] = append(ss.InterestedNodes[video], node)
}

func (ss *StreamingService) RemoveInterestedNode(video Video, node Node) {
	ss.Lock()
	defer ss.Unlock()
	for i, n := range ss.InterestedNodes[video] {
		if n == node {
			ss.InterestedNodes[video] = append(ss.InterestedNodes[video][:i], ss.InterestedNodes[video][i+1:]...)
			break
		}
	}
}
