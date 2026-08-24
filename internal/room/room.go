package room

import (
	"sync"

	"github.com/Nery93/chat-server/internal/client"
)

// Room significa que vamos cria uma sala
// para que os clientes possam se conectar e trocar mensagens
type Room struct {
	mu sync.Mutex
	Clients map[*client.Client]bool
}

func (r *Room) EntrarNaSala(c *client.Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Clients[c] = true
}

func (r *Room) SairDaSala(c *client.Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.Clients, c)
}

func (r *Room) Broadcast(message []byte) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for c := range r.Clients {
		select {
		case c.Send <- message:
		default:
			close(c.Send)
			delete(r.Clients, c)
		}
	}
}

func (r *Room) ClientCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.Clients)
}

func (r *Room) Usernames() []string {
	r.mu.Lock()
	defer r.mu.Unlock()

	usernames := make([]string, 0, len(r.Clients))
	for c := range r.Clients {
		usernames = append(usernames, c.Username)
	}

	return usernames
}

func NewRoom() *Room {
	return &Room{
		Clients: make(map[*client.Client]bool),
	}
}