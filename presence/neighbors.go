package presence

import (
	"fmt"
	"log"
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

func (nbl *NeighborList) PingNeighbors(cnf *config.AppConfigList, dvr *distancevectorrouting.DistanceVectorRouting, neighborsConnectionsMap *distancevectorrouting.NeighborsConnectionsMap) *distancevectorrouting.DistanceVectorRouting {

	msg := protobuf.Header{
		Type:      protobuf.RequestType_ROUTINGTABLE,
		Length:    0,
		Timestamp: int32(time.Now().UnixMilli()),
		Sender:    cnf.NodeIP.String(),
		Target:    "",
	}
	msg.Length = int32(proto.Size(&msg))

	var wg sync.WaitGroup
	results := make(chan *NeighborResult, len(nbl.content))

	for i := range nbl.content {
		neighbor := nbl.content[i]
		wg.Add(1)
		go func(neighbor *protobuf.Interface) {
			defer wg.Done()

			nb := distancevectorrouting.Interface{Interface: neighbor}
			fmt.Printf("Pinging %s\n", nb.ToString())

			// Start a new QUIC connection and stream if it doesn't exist already
			stream, conn, err := neighborsConnectionsMap.GetConnectionStream(nb.ToString())
			if err != nil {
				log.Printf("Error starting stream to %s: %v\n", nb.ToString(), err)
				return
			}

			msg.Target = nb.ToString()
			msg.Timestamp = int32(time.Now().UnixMilli())
			dvr.UpdateSource(distancevectorrouting.Interface{Interface: &protobuf.Interface{
				Ip:   cnf.NodeIP.String(),
				Port: cnf.NodeIP.Port,
			}})
			msg.Content = &protobuf.Header_DistanceVectorRouting{
				DistanceVectorRouting: dvr.Dvr,
			}

			fmt.Println("Sending ping to", nb.ToString())
			fmt.Println("Sender is", msg.Sender)

			data, err := proto.Marshal(&msg)
			if err != nil {
				log.Printf("Error marshaling ping: %v\n", err)
				return
			}

			if err := config.SendMessage(stream, data); err != nil {
				log.Printf("Error sending ping to %s: %v\n", nb.ToString(), err)
				return
			}

			responseData, err := config.ReceiveMessage(stream)
			if err != nil {
				log.Printf("Error receiving message from %s: %v\n", nb.ToString(), err)
				return
			}

			var response protobuf.Header
			if err := proto.Unmarshal(responseData, &response); err != nil {
				log.Printf("Error unmarshaling routing table from %s: %v\n", nb.ToString(), err)
				return
			}

			response.Sender = conn.RemoteAddr().String()
			fmt.Println("Received response from", response.Sender)

			if response.Type != protobuf.RequestType_ROUTINGTABLE {
				log.Printf("Error: received response is not a routing table from %s", nb.ToString())
				return
			}

			routingTable := response.GetDistanceVectorRouting()
			timeTook := time.UnixMilli(int64(response.Timestamp)).Sub(time.UnixMilli(int64(msg.Timestamp)))

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
	if len(nbl.content) > 1 {
		nbl.RemoveAll()
	}

  dvtList, delayList := processResults(results, *cnf)

	// Update the content of the routing table
	return distancevectorrouting.NewRouting(cnf.NodeIP, dvtList, delayList)
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
	n.content = append(n.content, neighbor)
}

func (n *NeighborList) RemoveAll() {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.content = make([]*protobuf.Interface, 0)
}

func (nl *NeighborList) RemoveNeighbor(neighbor *protobuf.Interface) {
	nl.mutex.Lock()
	defer nl.mutex.Unlock()
	for i, n := range nl.content {
		if n == neighbor {
			nl.content = append(nl.content[:i], nl.content[i+1:]...)
			break
		}
	}
}

func (n *NeighborList) HasNeighbor(neighbor *protobuf.Interface) bool {
	neighbor_i := distancevectorrouting.Interface{Interface: neighbor}
	for _, nb := range n.content {
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
		Mutex: sync.Mutex{},
		Dvr:   result.RoutingTable,
	}

	dvr.UpdateSource(distancevectorrouting.Interface{Interface: result.Neighbor})

	// single threaded
	*routingTables = append(*routingTables, dvr)

	*distances = append(*distances, distancevectorrouting.RequestRoutingDelay{
		Neighbor: distancevectorrouting.Interface{Interface: result.Neighbor},
		Delay:    int(result.Time),
	})
}

// appendSelf appends the routing table of the node itself.
func appendSelf(routingTables *[]distancevectorrouting.DistanceVectorRouting, distances *[]distancevectorrouting.RequestRoutingDelay, cnf config.AppConfigList) {
	myTable := distancevectorrouting.DistanceVectorRouting{
		Mutex: sync.Mutex{},
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
