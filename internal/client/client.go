package client

import (
	"encoding/json"
	"time"

	"github.com/Nery93/chat-server/internal/message"
	"github.com/gorilla/websocket"
)

const (
	pongWait   = 60 * time.Second
	writeWait  = 10 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type Client struct {
	Conn      *websocket.Conn
	Send      chan []byte
	Username  string
	Broadcast func(message []byte)
}

func (c *Client) WritePump() {

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := c.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}

	}
}

func (c *Client) ReadPump() {

	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	defer c.Conn.Close()
	for {
		_, raw, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		var msg message.Message
		if err := json.Unmarshal(raw, &msg); err != nil {
			continue
		}
		msg.Type = "chat"
		convert, err := json.Marshal(msg)
		if err != nil {
			continue
		}
		c.Broadcast(convert)
	}
}

func NewClient(conn *websocket.Conn, broadcast func(message []byte)) *Client {
	return &Client{
		Conn:      conn,
		Send:      make(chan []byte, 256),
		Broadcast: broadcast,
	}
}
