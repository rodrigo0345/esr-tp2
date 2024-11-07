package presence

import (
	"sync"

	"github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
)

type ClientList struct {
	mutex   sync.Mutex
	Content []*protobuf.Interface
}

func (cl *ClientList) Remove(client string) {
	cl.mutex.Lock()
	defer cl.mutex.Unlock()

	for i, c := range cl.Content {
		if c.String() == client {
			cl.Content = append(cl.Content[:i], cl.Content[i+1:]...)
			break
		}
	}
}

func (cl *ClientList) Add(client *protobuf.Interface) {
	cl.mutex.Lock()
	defer cl.mutex.Unlock()

	cl.Content = append(cl.Content, client)
}

func (cl *ClientList) Has(clientIp string) bool {
	clIP := config.ToInterface(clientIp)
	cl.mutex.Lock()
	defer cl.mutex.Unlock()

	for _, c := range cl.Content {
		if c.Ip == clIP.Ip && c.Port == clIP.Port {
			return true
		}
	}
	return false
}
