package api

// serveWs define our WebSocket endpoint
import (
	"fmt"
	"net/http"

	"github.com/SAIKAII/chatroom-backend/pkg/websocket"
)

func ServeWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &websocket.Client{
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	client.Read()
}
