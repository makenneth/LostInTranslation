package socket

import (
	"net/http"
	"log"
	"golang.org/x/net/websocket"
	// "encoding/json"
)

type SocketServer struct {
	path string
	clients map[string]*Client
	send chan *InterclientMessage
	addClient chan *Client
	removeClient chan *Client
}

func NewServer(path string) *SocketServer {
	clients := make(map[string]*Client, 0)
	addClient := make(chan *Client)
	removeClient := make(chan *Client)
	send := make(chan *InterclientMessage)

	return &SocketServer{path, clients, send, addClient, removeClient}
}


func (this *SocketServer) AddClient() chan<- *Client {
	return (chan<- *Client)(this.addClient)
}

func (this *SocketServer) RemoveClient() chan<- *Client {
	return (chan<- *Client)(this.removeClient)
}

func (this *SocketServer) Send() chan<- *InterclientMessage {
	return (chan<- *InterclientMessage)(this.send)
}

func (this *SocketServer) SendToClient(username string, message *Message) {
	if cl, ok := this.clients[username]; ok {
		cl.Write() <- message
	}
}

func (this *SocketServer) Listen() {
	log.Println("Socket Server started...")
	
	onConnect := func(ws *websocket.Conn){
		client := NewClient(ws, this)
		client.Listen()
		defer ws.Close()
	}

	http.Handle(this.path, websocket.Handler(onConnect))
	log.Println("Created handler")

	for {
		select {
		case cl := <- this.addClient:
			log.Printf("Adding new client %s", cl.username)
			this.clients[cl.username] = cl
			break;
		case cl := <- this.removeClient:
			log.Printf("Removing Client %s", cl.username)
			delete(this.clients, cl.username)
			break;
		case msg := <- this.send:
			for _, username := range msg.Dest {
				log.Printf("username dest: %s", username)
				log.Println("current Clients...", this.clients)
				if cl, ok := this.clients[username]; ok {
					log.Println("sending message to %s", cl.username)
					cl.Write() <- msg.Message
				}
			}
		}
	}
}
