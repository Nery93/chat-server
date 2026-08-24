package handler

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/Nery93/chat-server/internal/client"
	"github.com/Nery93/chat-server/internal/message"
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
		user := r.URL.Query().Get("user")

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

		var enterRoom message.Message
		enterRoom.Type = "join"
		enterRoom.User = user
		enterRoomJSON, err := json.Marshal(enterRoom)
		if err != nil {
			return
		}

		roomObj.Broadcast(enterRoomJSON)

		var userList message.Message
		userList.Type = "userlist"
		userList.Users = roomObj.Usernames()
		userListJSON, err := json.Marshal(userList)
		if err != nil {
			return
		}
		roomObj.Broadcast(userListJSON)

		go client.WritePump()
		client.ReadPump()

		var leaveRoom message.Message
		leaveRoom.Type = "leave"
		leaveRoom.User = user
		leaveRoomJSON, err := json.Marshal(leaveRoom)
		if err != nil {
			return
		}

		roomObj.Broadcast(leaveRoomJSON)
		roomObj.SairDaSala(client)
		
		var updatedUserList message.Message
		updatedUserList.Type = "userlist"
		updatedUserList.Users = roomObj.Usernames()
		updatedUserListJSON, err := json.Marshal(updatedUserList)
		if err != nil {
			return
		}
		roomObj.Broadcast(updatedUserListJSON)


	})

	mux.HandleFunc("GET /rooms", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()

		for roomName, roomObj := range rooms {
			if roomObj.ClientCount() == 0 {
				delete(rooms, roomName)
			}
		}

		roomList := make(map[string]int)

		for roomName, roomObj := range rooms {
			roomList[roomName] = roomObj.ClientCount()
		}

		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(roomList)
	})

	return mux
}
