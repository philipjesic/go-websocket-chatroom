package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}

type Client struct {
	conn *websocket.Conn
	send chan Message
	name string
}

type Hub struct {
	clients map[*Client]bool
	broadcast chan Message
	register chan *Client
	unregister chan *Client
	mu sync.Mutex
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections for demo purposes
	},
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[*Client]bool),
		broadcast: make(chan Message),
		register: make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <- h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("Client %s is connected", client.name)
		
		case client := <- h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("Client %s disconnected", client.name)
			}
			h.mu.Unlock()

		case message := <- h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client) 
				}
			}
			h.mu.Unlock()
		}
	}
}

func serveWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	// For demo: grab username from query param
	username := r.URL.Query().Get("username")
	if username == "" {
		username = "Anonymous"
	}

	client := &Client{
		conn: conn,
		send: make(chan Message, 256),
		name: username,
	}

	hub.register <- client

	// Read Messages from client, write to hub
	go func() {
		defer func() {
			hub.unregister <- client
			client.conn.Close()
		}()
		for {
			var message Message
			err := client.conn.ReadJSON(&message)
			if err != nil {
				log.Println("Read error: ", err)
				break
			}
			message.Username = client.name
			hub.broadcast <- message
		}
	}()

	// Write messages from hub to client
	go func(){
		for msg := range client.send {
			err := client.conn.WriteJSON(msg)
			if err != nil {
				log.Println("Write error: ", err)
				break
			}
		}
	}()
}


func main() {
	hub := NewHub()
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWebSocket(hub, w, r)
	})

	log.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
