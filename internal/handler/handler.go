package handler

import (
	"net/http"
	"sync"

	"github.com/Nery93/chat-server/internal/client"
	"github.com/Nery93/chat-server/internal/room"
	"github.com/gorilla/websocket"
)

var (
	mu    sync.Mutex
	rooms = make(map[string]*room.Room)
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewRouter() *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /ws/{sala_geral}", func(w http.ResponseWriter, r *http.Request) {

		sala := r.PathValue("sala_geral")
		user:= r.URL.Query().Get("user")

		if user == "" {
			user = "Anonymous"
		}

		mu.Lock()
		roomObj, exists := rooms[sala]

		if !exists {
			roomObj = room.NewRoom()
			rooms[sala] = roomObj
		}
		mu.Unlock()

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
			return
		}

		client := client.NewClient(conn, roomObj.Broadcast)
		client.Username = user
		roomObj.EntrarNaSala(client)
		go client.WritePump()
		client.ReadPump()
		roomObj.SairDaSala(client)
	})

	return mux
}
