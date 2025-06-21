package socket

import (
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/google/uuid"
	"github.com/hertz-contrib/websocket"
)

type Client struct {
	ID      uuid.UUID
	UserID  uuid.UUID
	GroupID uuid.UUID
	Conn    *websocket.Conn
	Send    chan []byte
}

func (c *Client) Read() {
	defer func() {
		SocketServer.UnRegister <- c
		c.Conn.Close()
	}()
	go c.handlePing()
	c.Conn.SetReadLimit(512)
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if strings.Contains(err.Error(), "1005") {
				// 正常关闭websocket
				hlog.Warnf("WebSocket connection closed normally: %v", err)
				break
			}
			if strings.Contains(err.Error(), "1006") {
				// 连接异常关闭
				hlog.Warnf("WebSocket connection closed abnormally: %v", err)
				break
			}
			hlog.Warnf("ReadMessage error: %v", err)
			break
		}
		if string(message) == `ping` {
			continue
		}
		if len(message) != 0 {
			hlog.Debug("send message to broadcast")
			SocketServer.Broadcast <- message
		}
	}
}

func (c *Client) Write() {
	defer func() {
		c.Conn.Close()
	}()
	for msg := range c.Send {
		// 为nil时已断开连接
		if c.Conn != nil {
			hlog.Debug("send data", msg)
			err := c.Conn.WriteMessage(websocket.BinaryMessage, msg)
			if err != nil {
				hlog.Error("WriteMessage error:", err)
				break
			}
		} else {
			break
		}
	}
}

func (c *Client) handlePing() {
	for {
		mt, message, err := c.Conn.ReadMessage()
		if err != nil {
			hlog.Error("read message error:", err)
			break
		}
		if mt == websocket.TextMessage && string(message) == `ping` {
			err = c.Conn.WriteMessage(websocket.TextMessage, []byte(`pong`))
			if err != nil {
				hlog.Warnf("Write pong error: %v", err)
				break
			}
		}
	}
}
