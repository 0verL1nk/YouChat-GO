package auth

import (
	"context"

	"core/biz/cerrors"
	service "core/biz/service/auth"
	"core/biz/utils"
	auth "core/hertz_gen/auth"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// Register .
// @router /register [POST]
func Register(ctx context.Context, c *app.RequestContext) {
	var err error
	var req auth.RegisterReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusBadRequest, err)
		return
	}

	resp := &auth.RegisterResp{}
	resp, err = service.NewRegisterService(ctx, c).Run(&req)
	if err != nil {
		if cerrors.Is(err, cerrors.BaseErrBadReq) {
			utils.SendErrResponse(ctx, c, consts.StatusBadRequest, err)
			return
		}
		utils.SendErrResponse(ctx, c, consts.StatusInternalServerError, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// Login .
// @router /login [POST]
func Login(ctx context.Context, c *app.RequestContext) {
	var err error
	var req auth.LoginReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusBadRequest, err)
		return
	}

	resp := &auth.LoginResp{}
	resp, err = service.NewLoginService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusInternalServerError, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// GetCaptcha .
// @router /auth/captcha [POST]
func GetCaptcha(ctx context.Context, c *app.RequestContext) {
	var err error
	var req auth.GetCaptchaReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusBadRequest, err)
		return
	}

	resp, err := service.NewGetCaptchaService(ctx, c).Run(&req)

	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusInternalServerError, err)
		return
	}
	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}
