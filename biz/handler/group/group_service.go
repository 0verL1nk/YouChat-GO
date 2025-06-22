package group

import (
	"context"

	service "core/biz/service/group"
	"core/biz/utils"
	group "core/hertz_gen/group"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// CreateGroup .
// @router /group/create [POST]
func CreateGroup(ctx context.Context, c *app.RequestContext) {
	var err error
	var req group.CreateGroupReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusBadRequest, err)
		return
	}

	resp := &group.CreateGroupResp{}
	resp, err = service.NewCreateGroupService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusInternalServerError, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// SearchGroup .
// @router /group/search [POST]
func SearchGroup(ctx context.Context, c *app.RequestContext) {
	var err error
	var req group.SearchGroupReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusBadRequest, err)
		return
	}

	resp := &group.SearchGroupResp{}
	resp, err = service.NewSearchGroupService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusInternalServerError, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// JoinGroup .
// @router /group/join [POST]
func JoinGroup(ctx context.Context, c *app.RequestContext) {
	var err error
	var req group.JoinGroupReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusBadRequest, err)
		return
	}

	resp, err := service.NewJoinGroupService(ctx, c).Run(&req)

	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusInternalServerError, err)
		return
	}
	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}
