package auth

import (
	"context"

	service "core/biz/service/auth"
	"core/biz/utils"
	auth "core/hertz_gen/auth"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// GetCheckCode .
// @router /checkCode [GET]
func GetCheckCode(ctx context.Context, c *app.RequestContext) {
	var err error
	var req auth.GetCheckCodeReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusBadRequest, err)
		return
	}

	resp := &auth.GetCheckCodeResp{}
	resp, err = service.NewGetCheckCodeService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusInternalServerError, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}
