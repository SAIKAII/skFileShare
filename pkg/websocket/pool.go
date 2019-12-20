package websocket

import "fmt"

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
}

// NewPool 新建一个Pool
func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

// Start 开始监听线上用户发送的消息，并把接收到的消息推送到每个用户
func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for client := range pool.Clients {
				client.Conn.WriteJSON(Message{
					Type: 1,
					Body: client.Conn.RemoteAddr().String() + "New User Joined...",
				})
			}
			pool.Clients[client] = true
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for client := range pool.Clients {
				client.Conn.WriteJSON(Message{
					Type: 1,
					Body: client.Conn.RemoteAddr().String() + "User Disconnected...",
				})
			}
		case message := <-pool.Broadcast:
			fmt.Println("Sending message to all clients in Pool")
			for client := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}
