package dashboard

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

// function is always ran in a go routine
// reads the clients send channel to detect if anything needs to be sent
func (c *Client) writePump() {
	for msg := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			c.conn.Close()
			break
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	clients = map[*Client]bool{}
	mu      sync.Mutex
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &Client{
		conn: conn,
		send: make(chan []byte, 256),
	}

	mu.Lock()
	clients[client] = true
	mu.Unlock()

	// goroutine infinite while loop checking the client's send channel
	// if anything in channel instantly send it
	go client.writePump()

	// goroutine to detect disconnects
	go func() {
		defer func() {
			mu.Lock()
			delete(clients, client)
			mu.Unlock()
			conn.Close()
		}()

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}()
}

func StartServer() {
	http.HandleFunc("/ws", wsHandler)
	http.ListenAndServe(":8090", nil)
}
