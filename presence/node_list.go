package presence

import (
	"sync"
	"time"
)

type NodeConnected struct {
	Node     string
	LastPing int64
}

type NodeList struct {
	mutex   sync.Mutex
	Content []NodeConnected
}

func (nl *NodeList) Has(node string) bool {
	nl.mutex.Lock()
	defer nl.mutex.Unlock()

	for _, n := range nl.Content {
		if n.Node == node {
			return true
		}
	}

	return false
}

func (nl *NodeList) Add(node string) {

	if nl.Has(node) {
		return
	}

	nl.mutex.Lock()
	defer nl.mutex.Unlock()

	if nl.Content == nil {
		nl.Content = make([]NodeConnected, 0)
	}

	lastPing := int64(time.Now().UnixMilli())

	nl.Content = append(nl.Content, NodeConnected{
		Node:     node,
		LastPing: lastPing,
	})
}

func (nl *NodeList) Remove(node string) {
	nl.mutex.Lock()
	defer nl.mutex.Unlock()

	for i, n := range nl.Content {
		if n.Node == node {
			nl.Content = append(nl.Content[:i], nl.Content[i+1:]...)
			break
		}
	}
}
