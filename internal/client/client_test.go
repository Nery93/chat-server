package client

import "testing"

func TestNewClient(t *testing.T) {
	broadcastCalls := 0
	broadcast := func(message []byte) { broadcastCalls++ }

	c := NewClient(nil, broadcast)

	if c.Send == nil {
		t.Fatal("Send channel not initialized")
	}
	if cap(c.Send) != 256 {
		t.Errorf("expected Send buffer of 256, got %d", cap(c.Send))
	}
	if c.Broadcast == nil {
		t.Fatal("Broadcast not set")
	}

	c.Broadcast([]byte("hello"))
	if broadcastCalls != 1 {
		t.Errorf("expected Broadcast to be called once, got %d", broadcastCalls)
	}
}
