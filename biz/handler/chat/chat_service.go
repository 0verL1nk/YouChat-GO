package chat

import (
	"context"

	service "core/biz/service/chat"
	"core/biz/utils"
	chat "core/hertz_gen/chat"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

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

	resp := &chat.ConnectChatWSResp{}
	resp, err = service.NewConnectChatWSService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusInternalServerError, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}
