package streaming

import (
	"log"
	"sync"
	"github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
	"google.golang.org/protobuf/proto"
)

type PubSub struct {
	mutex sync.Mutex
	// the string is the movie the client wants to watch
	clients        map[string]*RequestStream
	onGoingStreams []*RequestStream
}

func (ps *PubSub) notifyNewClient(client *RequestStream) {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()
	ps.clients[client.source] = client
}

func (ps *PubSub) notifyClientDisconnect(client *RequestStream) {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()
	delete(ps.clients, client.source)
}

func (ps *PubSub) notifyOnGoingStream(s *RequestStream) {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()
	ps.onGoingStreams = append(ps.onGoingStreams, s)

	go stream(s, ps)
}

func (ps *PubSub) notifyOnGoingStreamDisconnect(stream *RequestStream) {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()
	for i, s := range ps.onGoingStreams {
		if s == stream {
			ps.onGoingStreams = append(ps.onGoingStreams[:i], ps.onGoingStreams[i+1:]...)
			break
		}
	}
}

func (ps *PubSub) readStream(source string) *RequestStream {
	return ps.clients[source]
}

func stream(stream *RequestStream, pubSub *PubSub) {
  defer config.CloseStream(stream.stream)
	for {
		// get the video stream
		data, err := config.ReceiveMessage(stream.stream)
		if err != nil {
			log.Printf("Error receiving stream: %v", err)
			continue
		}

		chunk := &protobuf.Header{}
		err = proto.Unmarshal(data, chunk)

    if err != nil {
			log.Printf("Error unmarshalling chunk: %v", err)
			continue
		}

		// send the video stream to the interested clients
		for pubSub.clients[stream.source] != nil {
		}
	}

}
