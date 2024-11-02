package distancevectorrouting

import (
	"log"
	"sync"

	"github.com/quic-go/quic-go"
	"github.com/rodrigo0345/esr-tp2/config"
)

type ConnectionPool struct {
	Mutex       sync.Mutex
	Connections map[string]quic.Connection
}

func NewNeighborsConnectionsMap() *ConnectionPool {
	return &ConnectionPool{
		Mutex:       sync.Mutex{},
		Connections: make(map[string]quic.Connection),
	}
}

func (nm *ConnectionPool) Lock() {
	nm.Mutex.Lock()
}

func (nm *ConnectionPool) Unlock() {
	nm.Mutex.Unlock()
}

func (nm *ConnectionPool) Add(neighborIp string, conn quic.Connection) {
	nm.Lock()
	defer nm.Unlock()

	if nm.Connections == nil {
		nm.Connections = make(map[string]quic.Connection)
	}

	nm.Connections[neighborIp] = conn
}

func (nm *ConnectionPool) Remove(neighborIp string) {
	nm.Lock()
	defer nm.Unlock()

	if nm.Connections == nil {
		nm.Connections = make(map[string]quic.Connection)
	}

	delete(nm.Connections, neighborIp)
}

func (nm *ConnectionPool) CreateNewConnection(neighborIp string) {
	_, conn, err := config.StartConnStream(neighborIp)
	if err != nil {
		log.Printf("Error starting stream to %s: %v\n", neighborIp, err)
		return
	}
	// close the previous connection
	nm.Connections[neighborIp].CloseWithError(0, "")
	nm.Add(neighborIp, conn)
}

func (nm *ConnectionPool) GetConnectionStream(address string) (quic.Stream, quic.Connection, error) {

	var stream quic.Stream
	var conn quic.Connection
	var err error
	var ok bool = false

	if _, ok := nm.Connections[address]; ok {
		conn := nm.Connections[address]
		stream, _ = config.CreateStream(conn)
		// check if the connection is still alive
		ok = config.IsConnectionOpen(stream)
		if !ok {
			config.CloseStream(stream)
			config.CloseConnection(conn)
		}

	}

	if !ok {
		// create a new connection
		stream, conn, err = config.StartConnStream(address)
		if err != nil {
			nm.Remove(address)
			return nil, nil, err
		}

		// add the connection to the map
		nm.Add(address, conn)
	}

	// create a new connection
	return stream, conn, nil
}
