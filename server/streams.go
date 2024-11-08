package server

import "sync"

type Client struct {
	mutex            sync.Mutex
	PresenceNodeName string
	ClientIP         string
}
type Stream struct {
	// specifies what video and what client requested it
	Video string

	// list of clients that requested the video
	Clients []*Client
}

type VideoStreams struct {
	mutex sync.Mutex
	// map of video name to stream
	Streams []*Stream
}

func NewVideoStreams() *VideoStreams {
	return &VideoStreams{
		Streams: []*Stream{},
	}
}

// add a new stream to the list
func (vs *VideoStreams) AddStream(video string, client *Client) {
	vs.mutex.Lock()
	defer vs.mutex.Unlock()

	// check if the stream already exists
	for _, stream := range vs.Streams {
		if stream.Video == video {
			// add the client to the list
			stream.Clients = append(stream.Clients, client)
			return
		}
	}
	vs.Streams = append(vs.Streams, &Stream{
		Video: video,
		Clients: []*Client{
			client,
		},
	})
}

// remove a stream from the list
func (vs *VideoStreams) RemoveStream(video string, client Client) {
	vs.mutex.Lock()
	defer vs.mutex.Unlock()

	// check if the stream exists
	for i, stream := range vs.Streams {
		if stream.Video == video {
			// Remove the client from the list
			newClients := []*Client{}
			for _, c := range vs.Streams[i].Clients {
				if c.PresenceNodeName != client.PresenceNodeName {
					newClients = append(newClients, c)
				}
				// clean c
				c = nil
			}
			vs.Streams[i].Clients = newClients

			// If there are no more clients, remove the stream
			if len(vs.Streams[i].Clients) == 0 {
				// Remove the stream entirely if no clients are left
				vs.Streams = append(vs.Streams[:i], vs.Streams[i+1:]...)
				stream = nil
			}

			return // Assuming each video only appears once in the stream list
		}
	}
}

func (vs *VideoStreams) GetStreamClients(video string) []*Client {
	var clients []*Client
	vs.mutex.Lock()
	defer vs.mutex.Unlock()

	// check if the stream exists
	for _, stream := range vs.Streams {
		if stream.Video == video {
			clients = stream.Clients
			break
		}
	}
	return clients
}
