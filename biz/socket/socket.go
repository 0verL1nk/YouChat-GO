package socket

import (
	"github.com/hertz-contrib/websocket"
)

type MsgType string

const (
	MsgTypeText MsgType = "text"
)

type Message struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Type    string `json:"type"`
	Content string `json:"content"`
}
type Client struct {
	ID      string
	UserID  uint64
	GroupID string
	Conn    *websocket.Conn
	Send    chan []byte
}
