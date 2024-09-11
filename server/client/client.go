package client

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	name   string
	socket *websocket.Conn
}

func NewClient(ws *websocket.Conn, name string) *Client {
	return &Client{
		name:   name,
		socket: ws,
	}
}

func (c *Client) GetName() string {
	return c.name
}

func (c *Client) SendMessage(msg []byte) error {
	return c.socket.WriteMessage(websocket.TextMessage, msg)
}

func (c *Client) ReadMessage() (int, []byte, error) {
	return c.socket.ReadMessage()
}

func (c *Client) Close() error {
	c.socket.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	return c.socket.Close()
}

