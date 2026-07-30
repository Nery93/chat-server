package handler

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/Nery93/chat-server/internal/client"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewRouter() *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /ws/{sala}", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		client := client.NewClient(conn)
		go client.WritePump()
		client.ReadPump()
	})

	return mux
}
