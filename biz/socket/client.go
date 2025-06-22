package socket

import (
	"core/biz/cerrors"
	"strings"
	"sync"
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
	Closed  sync.Once
	Done    chan struct{} // 用于标记连接状态
	SendMu  sync.RWMutex
}

func NewClient(c *websocket.Conn, userID uuid.UUID) *Client {
	return &Client{
		ID:     userID, // 使用 userID 作为 Client 的 ID
		UserID: userID,
		Conn:   c,
		Send:   make(chan []byte, 10240),
		Done:   make(chan struct{}),
		SendMu: sync.RWMutex{},
	}
}

func (c *Client) Close() {
	c.Closed.Do(func() {
		close(c.Done)
		c.SendMu.Lock()
		if c.Send != nil {
			close(c.Send)
			c.Send = nil
		}
		c.SendMu.Unlock()
	})
}
func (c *Client) Read() {
	defer func() {
		SocketServer.UnRegister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(512)
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})
	for {
		mt, message, err := c.Conn.ReadMessage()
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
		if mt == websocket.TextMessage && string(message) == `ping` {
			err = c.Conn.WriteMessage(websocket.TextMessage, []byte(`pong`))
			if err != nil {
				hlog.Warnf("Write pong error: %v", err)
				break
			}
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
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second)) // 防止 Write 永久阻塞
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

func (c *Client) SafeSend(msg []byte) error {
	// 首先检查是否已关闭
	select {
	case <-c.Done:
		return cerrors.ErrClientClosed
	default:
	}

	c.SendMu.RLock()
	defer c.SendMu.RUnlock()

	// 再次检查 channel 是否为 nil
	if c.Send == nil {
		return cerrors.ErrClientClosed
	}

	// 尝试发送消息
	select {
	case c.Send <- msg:
		return nil
	case <-c.Done:
		return cerrors.ErrClientClosed
	default:
		return cerrors.ErrClientClosed // 或者创建一个新的错误类型表示channel已满
	}
}
