package model

type WSMessage struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Type    string `json:"type"`
	Code    int    `json:"code"`
	Content string `json:"content"`
}

type WSCode int

const (
	WSSuccess = 2000
	WSError   = 5000
)
