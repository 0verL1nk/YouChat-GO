package user

import (
	"context"

	service "core/biz/service/user"
	"core/biz/utils"
	user "core/hertz_gen/user"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// SearchUser .
// @router /user/search [POST]
func SearchUser(ctx context.Context, c *app.RequestContext) {
	var err error
	var req user.SearchUserReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusBadRequest, err)
		return
	}

	resp := &user.SearchUserResp{}
	resp, err = service.NewSearchUserService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusInternalServerError, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// GetUserInfo .
// @router /user/me [GET]
func GetUserInfo(ctx context.Context, c *app.RequestContext) {
	var err error
	var req user.GetUserInfoReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusBadRequest, err)
		return
	}

	resp := &user.GetUserInfoResp{}
	resp, err = service.NewGetUserInfoService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusInternalServerError, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// AddFriend .
// @router /user/add_friend [POST]
func AddFriend(ctx context.Context, c *app.RequestContext) {
	var err error
	var req user.AddFriendReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp, err := service.NewAddFriendService(ctx, c).Run(&req)

	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}
	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// GetUserContacts .
// @router /user/contacts [GET]
func GetUserContacts(ctx context.Context, c *app.RequestContext) {
	var err error
	var req user.GetUserContactsReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp, err := service.NewGetUserContactsService(ctx, c).Run(&req)

	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}
	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}
