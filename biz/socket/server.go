package socket

import (
	"core/biz/cerrors"
	"core/biz/chttp"
	"core/biz/dal/model"
	"core/biz/dal/query"
	"core/biz/service/group"
	mq_producer "core/biz/service/mq/producer"
	chat "core/hertz_gen/chat"
	"errors"
	"sync"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

type Server struct {
	// string键为userID
	Clients map[uuid.UUID]*Client
	// 互斥锁
	Mutex      *sync.Mutex
	Broadcast  chan []byte
	Register   chan *Client
	UnRegister chan *Client
	Capacity   int
	Size       int
}

func NewServer() *Server {
	return &Server{
		Clients:    make(map[uuid.UUID]*Client),
		Mutex:      &sync.Mutex{},
		Broadcast:  make(chan []byte, 10240),
		Register:   make(chan *Client, 100),
		UnRegister: make(chan *Client, 100),
	}
}

var SocketServer = NewServer()

func OkMsgResp(userID uuid.UUID) []byte {
	// msg := model.WSMessage{
	// 	Code:    chttp.MESSAGE_SUCCESS,
	// 	From:    "system",
	// 	To:      userID.String(),
	// 	Type:    model.MsgTypeText,
	// 	Content: "ok",
	// }
	// message, _ := json.Marshal(msg)
	msg := &chat.ChatMsg{
		Code:    chttp.MESSAGE_SUCCESS,
		From:    "system",
		To:      userID.String(),
		Type:    chat.ChatMsg_TEXT,
		Content: &chat.ChatMsg_Text{Text: "ok"},
	}
	message, _ := proto.Marshal(msg)
	hlog.Debug("ok message: ", message)
	return message
}

func ErrMsgResp(userID uuid.UUID, err error) []byte {
	content := chat.ChatMsg_Text{Text: err.Error()}
	msg := chat.ChatMsg{
		From:    "system",
		To:      userID.String(),
		Type:    chat.ChatMsg_TEXT,
		Content: &content,
	}
	message, _ := proto.Marshal(&msg)
	return message
}

func (s *Server) Start() {
	hlog.Info("socket server stated")
	for {
		select {
		case client := <-s.Register:
			s.Mutex.Lock()
			if _, ok := s.Clients[client.UserID]; ok {
				// 如果用户已存在，清理用户
				close(client.Send)
				delete(s.Clients, client.UserID)
			}
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
			hlog.Debug("broadcast message")
			msg := chat.ChatMsg{}
			if err := proto.Unmarshal(message, &msg); err != nil {
				hlog.Debugf("unmarshal message error: %v", err)
				break
			}

			if msg.To != "" {
				switch msg.Type {
				case chat.ChatMsg_TEXT:
					{
						cid := msg.Id
						msg.Id = uuid.NewString()
						// go SaveTextMsg(&msg)
						go mq_producer.HandlerWSMessage(&msg)
						// 将单聊视为两个人的房间
						go SendTextMsg(&msg)
						// 给client发送成功消息
						updateMsgID(msg.From, msg.Id, cid)
					}
				}
			} else {
				// TODO
				// 校验权限
				// 广播到全体
			}
			// s.Mutex.Lock()
			// 这个defer会导致死锁,因为for是无限的
			// defer s.Mutex.Unlock()
			// 此后处理需要锁的内容
			hlog.Debug("broadcast message end")
		}
	}
}

func updateMsgID(userID string, msgID string, cid string) (err error) {
	hlog.Debug("updateMsgID start")
	// parse uuid
	_userID, err := uuid.Parse(userID)
	if err != nil {
		return cerrors.ErrParseUUID
	}
	SocketServer.Mutex.Lock()
	userClient, ok := SocketServer.Clients[_userID]
	SocketServer.Mutex.Unlock()
	if !ok {
		return cerrors.ErrUserNotOnline
	}
	// 发送消息ID更新
	content := chat.ChatMsg_Text{Text: cid}
	message := chat.ChatMsg{
		Id:      msgID,
		From:    "msg_received",
		To:      userID,
		Type:    chat.ChatMsg_TEXT,
		Content: &content,
		Code:    chttp.MESSAGE_SUCCESS,
	}
	msg, _ := proto.Marshal(&message)
	userClient.Send <- msg
	hlog.Debug("updateMsgID end")
	return
}

func SendTextMsg(msg *chat.ChatMsg) (err error) {
	hlog.Debug("start send text message to group")

	// parse uuid
	fromID, err := uuid.Parse(msg.From)
	if err != nil {
		return cerrors.ErrParseUUID
	}
	toID, err := uuid.Parse(msg.To)
	if err != nil {
		return cerrors.ErrParseUUID
	}
	hlog.Debug("start get group user ids")
	ids, err := group.GetGroupUserIDs(toID, fromID)
	if err != nil {
		return err
	}
	hlog.Debug("get group user ids end")
	if len(ids) == 0 {
		hlog.Debug("end send text message to group")
		return nil
	}
	//
	hlog.Debug("start send text message to group")
	SocketServer.Mutex.Lock()
	// 发送消息给群组成员
	for _, id := range ids {
		if client, ok := SocketServer.Clients[id]; ok {
			message, _ := proto.Marshal(msg)
			select {
			case client.Send <- message: // 增加select防止通道阻塞
			default:
				hlog.Warn("client send channel full, dropping message", id)
			}
		}
	}
	SocketServer.Mutex.Unlock()
	hlog.Debug("send text message to group end")
	// 返回ok
	message := OkMsgResp(fromID)
	SocketServer.Mutex.Lock()
	client, ok := SocketServer.Clients[fromID]
	SocketServer.Mutex.Unlock()
	if ok {
		client.Send <- message
	}
	hlog.Debug("send text message to group end")
	return
}

var (
	ErrInvalidMsgType = errors.New("invalid message type")
	ErrParseUint      = errors.New("parse uint failed")
)

func SaveTextMsg(msg *chat.ChatMsg) (err error) {
	//parse uuid
	msgID, err := uuid.Parse(msg.Id)
	if err != nil {
		return cerrors.ErrParseUUID
	}
	// parse uuid
	// hlog.Debugf("fromID:%v toID:%v content:%s ", msg.From, msg.To, msg.Content)
	fromID, err := uuid.Parse(msg.From)
	if err != nil {
		return cerrors.ErrParseUUID
	}
	toID, err := uuid.Parse(msg.To)
	if err != nil {
		return cerrors.ErrParseUUID
	}
	// parse content
	content, ok := msg.Content.(*chat.ChatMsg_Text)
	if !ok {
		return ErrInvalidMsgType
	}
	err = query.Q.ChatMessage.Create(&model.ChatMessage{
		BaseModel: model.BaseModel{ID: msgID},
		MsgType:   model.MessageType(msg.Type),
		Content:   content.Text,
		FromId:    fromID,
		ToId:      toID,
	})
	if err != nil {
		return err
	}
	return
}
