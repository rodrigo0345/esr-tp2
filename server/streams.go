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
func (vs *VideoStreams) AddStream(video string, client *Client) *Stream {
	vs.mutex.Lock()
	defer vs.mutex.Unlock()

	// check if the stream already exists
	for _, stream := range vs.Streams {
		if stream.Video == video {
      // check if the client is already in the list
      for _, c := range stream.Clients {
        if c.PresenceNodeName == client.PresenceNodeName {
          return stream
        }
      }

			stream.Clients = append(stream.Clients, client)
			return stream
		}
	}

  s := Stream{
    Video: video,
    Clients: []*Client{
      client,
    },
  }
	vs.Streams = append(vs.Streams, &s)
  return &s
}

// remove a stream from the list
func (vs *VideoStreams) RemoveStream(video string, client Client) *Stream {
	vs.mutex.Lock()
	defer vs.mutex.Unlock()

  var s *Stream

	// check if the stream exists
	for i, stream := range vs.Streams {
		if stream.Video == video {
			// Remove the client from the list
      s = stream
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

      break
		}
	}
  return s
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
