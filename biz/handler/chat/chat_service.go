package chat

import (
	"context"
	"errors"
	"net/http"

	"core/biz/jwt"
	service "core/biz/service/chat"
	"core/biz/socket"
	"core/biz/utils"
	chat "core/hertz_gen/chat"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/websocket"
)

var upgrader = websocket.HertzUpgrader{CheckOrigin: func(ctx *app.RequestContext) bool {
	return true
}}

// ConnectChatWS .
// @router /ws/chat [GET]
func ConnectChatWS(ctx context.Context, c *app.RequestContext) {
	var err error
	var req chat.ConnectChatWSReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusBadRequest, err)
		return
	}
	tokenClaim, err := jwt.ParseToken(req.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errors.New("token invalid"))
	}
	if err = upgrader.Upgrade(c, func(c *websocket.Conn) {
		if c == nil {
			hlog.Error("websocket upgrade returned nil conn")
			return
		}
		client := socket.NewClient(c, tokenClaim.UserId)
		socket.SocketServer.Register <- client
		// 启动读写
		go client.Read()
		client.Write()
		// 这个用来阻塞线程
	}); err != nil {
		hlog.Debug("ws connect err: ", err)
		return
	}
}

// GetConversation .
// @router /chat/conversations/:groupID [GET]
func GetConversation(ctx context.Context, c *app.RequestContext) {
	var err error
	var req chat.GetConversationReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusBadRequest, err)
		return
	}

	resp, err := service.NewGetConversationService(ctx, c).Run(&req)

	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusInternalServerError, err)
		return
	}
	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}
