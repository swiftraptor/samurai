package internal

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Server struct {
	Clients map[string]*Client // id -> client
	Rooms   map[string]*MessageRoom
}

func MakeServer() Server {
	return Server{Clients: make(map[string]*Client, 256), Rooms: make(map[string]*MessageRoom, 64)}
}

func (server Server) wsHandler(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	id := uuid.New().String()

	client := &Client{Connection: conn, Sender: make(chan Message, 1), Id: id} // 1kb buffer?

	if _, ok := server.Clients[id]; ok {
		w.WriteHeader(400)
		return
	}

	server.Clients[id] = client
	log.Printf("User %s connected to server", id)
	go client.ReadStream(server)
	go client.WriteStream(server)

}

func (server Server) Start(host string, port int) {

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		server.wsHandler(w, r)
	})

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
