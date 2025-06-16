package socket

import (
	"strings"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/websocket"
)

type Client struct {
	ID      string
	UserID  uint64
	GroupID string
	Conn    *websocket.Conn
	Send    chan []byte
}

func (c *Client) Read() {
	defer func() {
		SocketServer.UnRegister <- c
		c.Conn.Close()
	}()
	for {
		mt, message, err := c.Conn.ReadMessage()
		if err != nil {
			if strings.Contains(err.Error(), "1005") {
				// 正常关闭websocket
				break
			}
			hlog.Warnf("ReadMessage error: %v", err)
			break
		}
		// 简单示例心跳处理
		if mt == websocket.TextMessage && string(message) == `ping` {
			err = c.Conn.WriteMessage(websocket.TextMessage, []byte(`pong`))
			if err != nil {
				hlog.Warnf("Write pong error: %v", err)
				break
			}
			continue
		}
		// 解码
		SocketServer.Broadcast <- message
	}
}

func (c *Client) Write() {
	defer func() {
		c.Conn.Close()
	}()
	for msg := range c.Send {
		// 为nil时已断开连接
		if c.Conn != nil {
			err := c.Conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				break
			}
		} else {
			break
		}
	}
}
