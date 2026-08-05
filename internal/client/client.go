package client


import (
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
	Send chan []byte
	Broadcast func(message []byte)
}

func (c *Client) WritePump(){

	for {
		message, ok := <-c.Send
		if !ok {
			c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		err := c.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			return
		}
	}
}

func (c *Client) ReadPump() {
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		c.Broadcast(message)
	}
}


func NewClient(conn *websocket.Conn, broadcast func(message []byte)) *Client {
	return &Client{
		Conn:      conn,
		Send:      make(chan []byte, 256),
		Broadcast: broadcast,
	}
}