package socket

import (
	"core/biz/dal/model"
	"core/biz/dal/query"
	"encoding/json"
	"errors"
	"strconv"
	"sync"

	"core/biz/service/group"
)

type Server struct {
	// string键为userID
	Clients map[uint64]*Client
	// 互斥锁
	Mutex      *sync.Mutex
	Broadcast  chan []byte
	Register   chan *Client
	UnRegister chan *Client
}

func NewServer() *Server {
	return &Server{
		Clients:    make(map[uint64]*Client),
		Mutex:      &sync.Mutex{},
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		UnRegister: make(chan *Client),
	}
}

var SocketServer = NewServer()

func OkMsgResp(userID uint64) []byte {
	msg := Message{
		From:    "system",
		To:      strconv.FormatUint(userID, 10),
		Type:    string(MsgTypeText),
		Content: "ok",
	}
	message, _ := json.Marshal(msg)
	return message
}

func (s *Server) Start() {
	for {
		select {
		case client := <-s.Register:
			s.Mutex.Lock()
			s.Clients[client.UserID] = client
			message := OkMsgResp(client.UserID)
			client.Send <- message
			s.Mutex.Unlock()
		case client := <-s.UnRegister:
			s.Mutex.Lock()
			if _, ok := s.Clients[client.UserID]; ok {
				close(client.Send)
				delete(s.Clients, client.UserID)
			}
			s.Mutex.Unlock()
		case message := <-s.Broadcast:
			s.Mutex.Lock()
			msg := Message{}
			json.Unmarshal(message, &msg)
			if msg.To != "" {
				switch msg.Type {
				case string(MsgTypeText):
					{
						SaveTextMsg(msg)
						// 将单聊视为两个人的房间
						SendTextMsg(msg)
					}
				}
			}
			s.Mutex.Unlock()
		}
	}
}

func SendTextMsg(msg Message) (err error) {
	SocketServer.Mutex.Lock()
	defer SocketServer.Mutex.Unlock()
	var toID uint64
	var fromID uint64
	toID, err = strconv.ParseUint(msg.To, 10, 64)
	if err != nil {
		return ErrParseUint
	}
	fromID, err = strconv.ParseUint(msg.From, 10, 64)
	if err != nil {
		return ErrParseUint
	}
	ids, err := group.GetGroupUserIDs(toID, fromID)
	if err != nil {
		return err
	}
	if len(ids) == 0 {
		return nil
	}
	for _, id := range ids {
		if client, ok := SocketServer.Clients[id]; ok {
			message, _ := json.Marshal(msg)
			client.Send <- message
		}
	}
	// 返回ok
	message := OkMsgResp(fromID)
	if client, ok := SocketServer.Clients[fromID]; ok {
		client.Send <- message
	}
	return
}

var (
	ErrInvalidMsgType = errors.New("invalid message type")
	ErrParseUint      = errors.New("parse uint failed")
)

func SaveTextMsg(msg Message) (err error) {
	msgType, ok := model.Str2MsgType[msg.Type]
	if !ok {
		return ErrInvalidMsgType
	}
	fromID, err := strconv.ParseUint(msg.From, 10, 64)
	if err != nil {
		return ErrParseUint
	}
	toID, err := strconv.ParseUint(msg.To, 10, 64)
	if err != nil {
		return ErrParseUint
	}

	err = query.Q.ChatMessage.Create(&model.ChatMessage{
		MsgType: msgType,
		Content: msg.Content,
		FromId:  fromID,
		ToId:    toID,
	})
	if err != nil {
		return err
	}
	return
}
