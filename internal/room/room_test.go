package room

import (
	"testing"

	"github.com/Nery93/chat-server/internal/client"
)

func TestRoom_AddClient(t *testing.T) {
	room := NewRoom()

	c := &client.Client{
		Send: make(chan []byte, 1),
	}

	room.EntrarNaSala(c)

	if len(room.Clients) != 1 {
		t.Errorf("Expected 1 client in the room, got %d", len(room.Clients))
	}
}

func TestRoom_RemoveClient(t *testing.T) {
	room := NewRoom()

	c := &client.Client{
		Send: make(chan []byte, 1),
	}

	room.EntrarNaSala(c)
	room.SairDaSala(c)

	if len(room.Clients) != 0 {
		t.Errorf("Expected 0 clients in the room, got %d", len(room.Clients))
	}
}

func TestRoom_Broadcast(t *testing.T) {
	room := NewRoom()

	c1 := &client.Client{
		Send: make(chan []byte, 1),
	}
	c2 := &client.Client{
		Send: make(chan []byte, 1),
	}

	room.EntrarNaSala(c1)
	room.EntrarNaSala(c2)

	message := []byte("Hello, World!")
	room.Broadcast(message)

	select {
	case msg := <-c1.Send:
		if string(msg) != string(message) {
			t.Errorf("Expected message '%s', got '%s'", message, msg)
		}
	default:
		t.Error("Expected a message for client 1, but got none")
	}

	select {
	case msg := <-c2.Send:
		if string(msg) != string(message) {
			t.Errorf("Expected message '%s', got '%s'", message, msg)
		}
	default:
		t.Error("Expected a message for client 2, but got none")
	}
}