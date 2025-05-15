package socket

import "github.com/hertz-contrib/websocket"

func (c *Client) Read() {
	defer func() {
		SocketServer.UnRegister <- c
		c.Conn.Close()
	}()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		SocketServer.Broadcast <- message
	}
}

func (c *Client) Write() {
	defer func() {
		c.Conn.Close()
	}()
	for msg := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			break
		}
	}
}
