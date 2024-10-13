package support

import (
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

func Support() {
  http.HandleFunc("/ws", handleWebSocket)
	http.Handle("/", http.FileServer(http.Dir(".")))

	log.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	udpAddr, err := net.ResolveUDPAddr("udp", ":5004")
	if err != nil {
		log.Println(err)
		return
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Println(err)
		return
	}
	defer udpConn.Close()

	buffer := make([]byte, 4096)
	for {
		n, _, err := udpConn.ReadFromUDP(buffer)
		if err != nil {
			log.Println(err)
			return
		}

		err = conn.WriteMessage(websocket.BinaryMessage, buffer[:n])
		if err != nil {
			log.Println(err)
			return
		}
	}
}
