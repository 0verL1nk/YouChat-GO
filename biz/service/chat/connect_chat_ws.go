package chat

import (
	"context"
	"errors"

	"core/biz/jwt"
	"core/biz/socket"
	chat "core/hertz_gen/chat"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/google/uuid"
	"github.com/hertz-contrib/websocket"
)

type ConnectChatWSService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

var upgrader = websocket.HertzUpgrader{}

func NewConnectChatWSService(Context context.Context, RequestContext *app.RequestContext) *ConnectChatWSService {
	return &ConnectChatWSService{RequestContext: RequestContext, Context: Context}
}

func (h *ConnectChatWSService) Run(req *chat.ConnectChatWSReq) (resp *chat.ConnectChatWSResp, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code
	id := uuid.New().String()
	tokenClaim, exist := h.RequestContext.Get("token")
	if !exist {
		return &chat.ConnectChatWSResp{}, errors.New("token is missing")
	}
	claims, ok := tokenClaim.(*jwt.TokenClaims)
	if !ok {
		return &chat.ConnectChatWSResp{}, errors.New("token is invalid")
	}
	err = upgrader.Upgrade(h.RequestContext, func(c *websocket.Conn) {
		client := &socket.Client{
			ID:     id,
			UserID: claims.UserId,
			Conn:   c,
			Send:   make(chan []byte),
		}
		socket.SocketServer.Register <- client
		go client.Read()
		go client.Write()
	})
	return
}
