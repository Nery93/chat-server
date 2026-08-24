package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Nery93/chat-server/internal/message"
	"github.com/gorilla/websocket"
)

func dial(t *testing.T, wsURL, room, user string) *websocket.Conn {
	t.Helper()
	url := wsURL + "/ws/" + room + "?user=" + user
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("dial failed: %v", err)
	}
	t.Cleanup(func() { conn.Close() })
	return conn
}

func readMessage(t *testing.T, conn *websocket.Conn) message.Message {
	t.Helper()
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	_, raw, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}
	var msg message.Message
	if err := json.Unmarshal(raw, &msg); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	return msg
}

func TestJoinAndBroadcast(t *testing.T) {
	server := httptest.NewServer(NewRouter())
	defer server.Close()
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	connA := dial(t, wsURL, "testroom", "Ana")
	if msg := readMessage(t, connA); msg.Type != "join" || msg.User != "Ana" {
		t.Fatalf("expected Ana's own join, got %+v", msg)
	}
	if msg := readMessage(t, connA); msg.Type != "userlist" {
		t.Fatalf("expected userlist after Ana's own join, got %+v", msg)
	}

	connB := dial(t, wsURL, "testroom", "Bruno")

	if msg := readMessage(t, connA); msg.Type != "join" || msg.User != "Bruno" {
		t.Fatalf("expected Bruno's join for Ana, got %+v", msg)
	}
	if msg := readMessage(t, connA); msg.Type != "userlist" {
		t.Fatalf("expected userlist for Ana after Bruno joined, got %+v", msg)
	}
	if msg := readMessage(t, connB); msg.Type != "join" || msg.User != "Bruno" {
		t.Fatalf("expected Bruno's own join, got %+v", msg)
	}
	if msg := readMessage(t, connB); msg.Type != "userlist" {
		t.Fatalf("expected userlist for Bruno, got %+v", msg)
	}

	if err := connA.WriteMessage(websocket.TextMessage, []byte(`{"type":"chat","text":"ola"}`)); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	if msg := readMessage(t, connA); msg.Type != "chat" || msg.User != "Ana" || msg.Text != "ola" {
		t.Fatalf("expected Ana's own echoed chat, got %+v", msg)
	}
	if msg := readMessage(t, connB); msg.Type != "chat" || msg.User != "Ana" || msg.Text != "ola" {
		t.Fatalf("expected Bruno to receive Ana's chat, got %+v", msg)
	}
}

func TestUserIdentityForced(t *testing.T) {
	server := httptest.NewServer(NewRouter())
	defer server.Close()
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	connA := dial(t, wsURL, "identitytest", "Ana")
	readMessage(t, connA) // Ana's own join
	readMessage(t, connA) // userlist after Ana's own join

	connB := dial(t, wsURL, "identitytest", "Bruno")
	readMessage(t, connA) // Bruno's join
	readMessage(t, connA) // userlist after Bruno joined
	readMessage(t, connB) // Bruno's own join
	readMessage(t, connB) // userlist after Bruno's own join

	spoof := `{"type":"system","user":"Bruno","text":"mensagem falsa"}`
	if err := connA.WriteMessage(websocket.TextMessage, []byte(spoof)); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	msg := readMessage(t, connB)
	if msg.Type != "chat" {
		t.Errorf("expected type forced to chat, got %q", msg.Type)
	}
	if msg.User != "Ana" {
		t.Errorf("expected user forced to Ana, got %q", msg.User)
	}
}

func TestGetRooms(t *testing.T) {
	server := httptest.NewServer(NewRouter())
	defer server.Close()
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	dial(t, wsURL, "roomslisttest", "Ana")

	resp, err := http.Get(server.URL + "/rooms")
	if err != nil {
		t.Fatalf("GET /rooms failed: %v", err)
	}
	defer resp.Body.Close()

	var rooms map[string]int
	if err := json.NewDecoder(resp.Body).Decode(&rooms); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if rooms["roomslisttest"] != 1 {
		t.Errorf("expected 1 client in roomslisttest, got %d", rooms["roomslisttest"])
	}
}
