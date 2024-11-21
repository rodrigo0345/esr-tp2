package presence

import (
	"fmt"
	"log"
	"math"
	_ "net"
	"sync"
	"time"

	"github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
	distancevectorrouting "github.com/rodrigo0345/esr-tp2/presence/distance_vector_routing"
	"google.golang.org/protobuf/proto"
)

type NeighborResult struct {
	Neighbor     *protobuf.Interface
	Time         time.Duration
	RoutingTable *protobuf.DistanceVectorRouting
}

func (nbl *NeighborList) PingNeighbors(logger *config.Logger, cnf *config.AppConfigList, dvr *distancevectorrouting.DistanceVectorRouting, neighborsConnectionsMap *distancevectorrouting.ConnectionPool) *distancevectorrouting.DistanceVectorRouting {

	msg := protobuf.Header{
		Type:      protobuf.RequestType_ROUTINGTABLE,
		Length:    0,
		Timestamp: time.Now().UnixMilli(),
		Sender:    cnf.NodeIP.String(),
		Target:    "",
	}
	msg.Length = int32(proto.Size(&msg))

	dvr.Dvr.Source = &protobuf.Interface{
		Ip:   cnf.NodeIP.Ip,
		Port: cnf.NodeIP.Port,
	}

	var wg sync.WaitGroup
	results := make(chan *NeighborResult, len(nbl.Content))

	for i := range nbl.Content {

		neighbor := nbl.Content[i]

		wg.Add(1)
		go func(neighbor *protobuf.Interface) {
			defer wg.Done()

			nb := distancevectorrouting.Interface{Interface: neighbor}

			// Start a new QUIC connection and stream if it doesn't exist already
			stream, conn, err := neighborsConnectionsMap.GetConnectionStream(nb.ToString())
			if err != nil {
        // logger.Error(fmt.Sprintf("Error starting stream to %s: %v", nb.ToString(), err))
				r, err := markNeighborAsDisconnected(dvr, nb)
				if err != nil {
					return
				}
				results <- r
				return
			}

			msg.Target = nb.ToString()
			msg.Content = &protobuf.Header_DistanceVectorRouting{
				DistanceVectorRouting: dvr.Dvr,
			}

			data, err := proto.Marshal(&msg)

			if err != nil {
        logger.Error(fmt.Sprintf("Error marshaling ping: %v", err))
				return
			}

			if err := config.SendMessage(stream, data); err != nil {
        logger.Error(fmt.Sprintf("Error sending ping to %s: %v", nb.ToString(), err))
				r, err := markNeighborAsDisconnected(dvr, nb)
				if err != nil {
					return
				}
				results <- r
				return
			}

			responseData, err := config.ReceiveMessage(stream)
			if err != nil {
        logger.Error(fmt.Sprintf("Error receiving message from %s: %v", nb.ToString(), err))
				r, err := markNeighborAsDisconnected(dvr, nb)
				if err != nil {
					return
				}
				results <- r
				return
			}

			var response protobuf.Header
			if err := proto.Unmarshal(responseData, &response); err != nil {
        logger.Error(fmt.Sprintf("Error unmarshaling routing table from %s: %v", nb.ToString(), err))
				return
			}

			response.Sender = conn.RemoteAddr().String()

			if response.Type != protobuf.RequestType_ROUTINGTABLE {
				log.Printf("Error: received response is not a routing table from %s", nb.ToString())
				return
			}

			routingTable := response.GetDistanceVectorRouting()
			timeTook := time.UnixMilli(response.Timestamp).Sub(time.UnixMilli(msg.Timestamp))

			// Log successful response for debugging
			// log.Printf("Received routing table from %s, time took: %v", nb.ToString(), timeTook)
			in := config.ToInterface(response.Sender)

			results <- &NeighborResult{
				Neighbor:     in,
				Time:         timeTook,
				RoutingTable: routingTable,
			}
		}(neighbor)
	}

	wg.Wait()
	close(results)

	// update the neighbor list
	// if len(nbl.content) > 1 {
	// 	nbl.RemoveAll()
	// }

	dvtList, delayList := processResults(results, *cnf)

	// Update the content of the routing table
	return distancevectorrouting.NewRouting(cnf.NodeName, cnf.NodeIP, dvtList, delayList)
}

func markNeighborAsDisconnected(dvr *distancevectorrouting.DistanceVectorRouting, nb distancevectorrouting.Interface) (*NeighborResult, error) {
	// search the dvr values for the neighborIp
	neighborName, err := dvr.GetName(&protobuf.Interface{Ip: nb.Ip, Port: nb.Port})
	neighborIp, err := dvr.GetNextHop(neighborName)
	if err != nil {
		return nil, err
	}

	nbIP := neighborIp.NextNode

	table := &protobuf.DistanceVectorRouting{
		Entries: make(map[string]*protobuf.NextHop),
		Source:  nbIP,
	}

	table.Entries[neighborName] = &protobuf.NextHop{NextNode: nbIP, Distance: math.MaxInt64}

	return &NeighborResult{
		Neighbor:     nbIP,
		Time:         time.Duration(math.MaxInt64),
		RoutingTable: table,
	}, nil
}

func (n *NeighborList) AddNeighbor(neighbor *protobuf.Interface, cnf *config.AppConfigList) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if n.HasNeighbor(neighbor) {
		return
	}

	if neighbor.Ip == cnf.NodeIP.Ip && neighbor.Port == cnf.NodeIP.Port {
		return
	}
	n.Content = append(n.Content, neighbor)
}

func (n *NeighborList) RemoveAll() {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.Content = make([]*protobuf.Interface, 0)
}

func (nl *NeighborList) RemoveNeighbor(neighbor *protobuf.Interface) {
	nl.mutex.Lock()
	defer nl.mutex.Unlock()
	for i, n := range nl.Content {
		if n == neighbor {
			nl.Content = append(nl.Content[:i], nl.Content[i+1:]...)
			break
		}
	}
}

func (n *NeighborList) HasNeighbor(neighbor *protobuf.Interface) bool {
	neighbor_i := distancevectorrouting.Interface{Interface: neighbor}
	for _, nb := range n.Content {
		nb_i := distancevectorrouting.Interface{Interface: nb}
		if nb_i.ToString() == neighbor_i.ToString() {
			return true
		}
	}
	return false
}

func processResults(results chan *NeighborResult, cnf config.AppConfigList) ([]distancevectorrouting.DistanceVectorRouting, []distancevectorrouting.RequestRoutingDelay) {
	var (
		routingTables []distancevectorrouting.DistanceVectorRouting
		distances     []distancevectorrouting.RequestRoutingDelay
	)

	for result := range results {
		addResult(&routingTables, &distances, *result, cnf)
	}

	// Append the node's own routing information
	appendSelf(&routingTables, &distances, cnf)

	return routingTables, distances
}

// addResult adds a single result to the routing tables and distances.
func addResult(routingTables *[]distancevectorrouting.DistanceVectorRouting, distances *[]distancevectorrouting.RequestRoutingDelay, result NeighborResult, cnf config.AppConfigList) {
	dvr := distancevectorrouting.DistanceVectorRouting{
		Mutex: sync.RWMutex{},
		Dvr:   result.RoutingTable,
	}

	dvr.UpdateSource(distancevectorrouting.Interface{Interface: result.Neighbor})

	// single threaded
	*routingTables = append(*routingTables, dvr)

	*distances = append(*distances, distancevectorrouting.RequestRoutingDelay{
		Neighbor: distancevectorrouting.Interface{Interface: result.Neighbor},
		Delay:    int64(result.Time),
	})
}

// appendSelf appends the routing table of the node itself.
func appendSelf(routingTables *[]distancevectorrouting.DistanceVectorRouting, distances *[]distancevectorrouting.RequestRoutingDelay, cnf config.AppConfigList) {
	myTable := distancevectorrouting.DistanceVectorRouting{
		Mutex: sync.RWMutex{},
		Dvr: &protobuf.DistanceVectorRouting{
			Source: cnf.NodeIP,
			Entries: map[string]*protobuf.NextHop{
				cnf.NodeName: {NextNode: cnf.NodeIP, Distance: 0},
			},
		},
	}

	// single threaded
	*routingTables = append(*routingTables, myTable)
	*distances = append(*distances, distancevectorrouting.RequestRoutingDelay{
		Neighbor: distancevectorrouting.Interface{Interface: cnf.NodeIP},
		Delay:    0,
	})
}
