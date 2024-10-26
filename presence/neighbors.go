package presence

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
	distancevectorrouting "github.com/rodrigo0345/esr-tp2/presence/distance_vector_routing"
	"google.golang.org/protobuf/proto"
)

func (nbl *NeighborList) PingNeighbors(cnf *config.AppConfigList, dvr *distancevectorrouting.DistanceVectorRouting) {
	fmt.Println("Synchronizing neighbors...")

	msg := protobuf.Header{
		Type:      protobuf.RequestType_ROUTINGTABLE,
		Length:    0,
		Timestamp: int32(time.Now().UnixMilli()),
		Sender:    cnf.NodeIP.String(),
		Target:    "",
		Content: &protobuf.Header_DistanceVectorRouting{
			DistanceVectorRouting: dvr.DistanceVectorRouting,
		},
	}
	msg.Length = int32(proto.Size(&msg))

	type NeighborResult struct {
		Neighbor     *protobuf.Interface
		Time         time.Duration
		RoutingTable *protobuf.DistanceVectorRouting
	}

	var wg sync.WaitGroup
	results := make(chan *NeighborResult, len(nbl.content))

	for i := range nbl.content {
		neighbor := nbl.content[i]
		wg.Add(1)
		go func(neighbor *protobuf.Interface) {
			defer wg.Done()

			nb := distancevectorrouting.Interface{Interface: neighbor}
			fmt.Printf("Pinging %s\n", nb.ToString())

			msg.Target = nb.ToString()
			msg.Timestamp = int32(time.Now().UnixMilli())

			data, err := proto.Marshal(&msg)
			if err != nil {
				log.Printf("Error marshaling ping: %v\n", err)
				return
			}

			// Start a new QUIC stream
			stream, err := config.StartStream(nb.ToString())
			if err != nil {
				log.Printf("Error starting stream to %s: %v\n", nb.ToString(), err)
				return
			}
			defer config.CloseStream(stream)

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

			if response.Type != protobuf.RequestType_ROUTINGTABLE {
				log.Printf("Error: received response is not a routing table from %s", nb.ToString())
				return
			}

			routingTable := response.GetDistanceVectorRouting()
			timeTook := time.UnixMilli(int64(response.Timestamp)).Sub(time.UnixMilli(int64(msg.Timestamp)))

			// Log successful response for debugging
			log.Printf("Received routing table from %s, time took: %v", nb.ToString(), timeTook)

			results <- &NeighborResult{
				Neighbor:     neighbor,
				Time:         timeTook,
				RoutingTable: routingTable,
			}
		}(neighbor)
	}

	wg.Wait()
	close(results)

	// Process results from the channel
	routingTables := make([]distancevectorrouting.DistanceVectorRouting, 0)
	distances := make([]distancevectorrouting.RequestRoutingDelay, 0)
	for result := range results {
		nbl.AddNeighbor(result.Neighbor)
		routingTables = append(routingTables, distancevectorrouting.DistanceVectorRouting{result.RoutingTable})

		distances = append(distances, distancevectorrouting.RequestRoutingDelay{
			Neighbor: distancevectorrouting.Interface{Interface: result.Neighbor},
			Delay:    int(result.Time),
		})
	}

	// Also append itself to the routing table
	tmp := map[string]*protobuf.NextHop{
		distancevectorrouting.Interface{Interface: cnf.NodeIP}.ToString(): {
			NextNode: cnf.NodeIP,
			Distance: 0,
		},
	}
	myTable := distancevectorrouting.DistanceVectorRouting{
		DistanceVectorRouting: &protobuf.DistanceVectorRouting{
			Source:  cnf.NodeIP,
			Entries: tmp,
		},
	}
	routingTables = append(routingTables, myTable)
	distances = append(distances, distancevectorrouting.RequestRoutingDelay{
		Neighbor: distancevectorrouting.Interface{Interface: cnf.NodeIP},
		Delay:    0,
	})

	// Update the content of the routing table
	*dvr = *distancevectorrouting.NewRouting(cnf.NodeIP, routingTables, distances)
}

func (n *NeighborList) AddNeighbor(neighbor *protobuf.Interface) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	// check if the neighbor is already in the list
	if n.HasNeighbor(neighbor) {
		return
	}
	n.content = append(n.content, neighbor)
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

