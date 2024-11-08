package presence

import (
	"fmt"
	"sync"
	"time"

	"github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
)

type ClientList struct {
	mutex    sync.Mutex
	Content  []*protobuf.Interface
	LastPing map[string]int
}

func (cl *ClientList) Remove(client string) {
	cl.mutex.Lock()
	defer cl.mutex.Unlock()

	for i, c := range cl.Content {
		if fmt.Sprintf("%s:%d", c.Ip, c.Port) == client {
			cl.Content = append(cl.Content[:i], cl.Content[i+1:]...)
			break
		}
	}

	// remove from the last ping map
	delete(cl.LastPing, client)
}

func (cl *ClientList) Add(client *protobuf.Interface) {
	if cl.LastPing == nil {
		cl.LastPing = make(map[string]int)
	}

	cl.mutex.Lock()
	defer cl.mutex.Unlock()

	// check if the client is already in the list
	for _, c := range cl.Content {
		if c.String() == client.String() {
			return
		}
	}

	cl.Content = append(cl.Content, client)
	cl.LastPing[fmt.Sprintf("%s:%d", client.Ip, client.Port)] = int(time.Now().Unix())
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

func (cl *ClientList) clientKey(ip string, port int) string {
	return fmt.Sprintf("%s:%d", ip, port)
}

func (cl *ClientList) ResetPing(client string) {
	cl.mutex.Lock()
	defer cl.mutex.Unlock()

	cl.LastPing[client] = int(time.Now().Unix())
}

func (cl *ClientList) GetPing(clientIp string, clientPort int) (int, bool) {
	cl.mutex.Lock()
	defer cl.mutex.Unlock()

	key := cl.clientKey(clientIp, clientPort)
	ping, exists := cl.LastPing[key]
	return ping, exists
}

// clients that don't ping within a specified interval
func (cl *ClientList) GetDeadClients(interval int) []*protobuf.Interface {
	var deadClients []*protobuf.Interface

	for _, client := range cl.Content {
		key := cl.clientKey(client.Ip, int(client.Port))
		if lastPing, exists := cl.LastPing[key]; exists {
			secondsPassed := int(time.Now().Unix()) - lastPing
			if secondsPassed > interval {
				deadClients = append(deadClients, client)
			}
		} else {
			// Optionally, consider clients without LastPing as "dead" immediately
			deadClients = append(deadClients, client)
		}
	}

	return deadClients
}

func (cl *ClientList) RemoveDeadClients(interval int) {

	for _, client := range cl.GetDeadClients(interval) {
		cl.Remove(cl.clientKey(client.Ip, int(client.Port)))
	}
}
