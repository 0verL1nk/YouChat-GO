package group

import (
	"context"

	"core/biz/cerrors"
	"core/biz/dal/model"
	"core/biz/dal/query"
	"core/biz/service/user"
	"core/biz/utils"
	group "core/hertz_gen/group"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/google/uuid"
)

type CreateGroupService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewCreateGroupService(Context context.Context, RequestContext *app.RequestContext) *CreateGroupService {
	return &CreateGroupService{RequestContext: RequestContext, Context: Context}
}

func (h *CreateGroupService) Run(req *group.CreateGroupReq) (resp *group.CreateGroupResp, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code
	// 获取userID
	userToken, err := utils.GetTokenFromMiddleware(h.RequestContext)
	if err != nil {
		hlog.Debugf("err get token from middleware: %v", err)
		return nil, err
	}
	hlog.Debug("create group userID:", userToken.UserId)
	// 查看用户是否存在
	if err := user.CheckUserExist(h.Context, userToken.UserId); err != nil {
		hlog.Debugf("err check user exist:%v", err)
		return nil, err
	}
	// 检查名字是否重复
	num, err := query.Group.WithContext(h.Context).Where(query.Group.GroupName.Eq(req.Name)).Count()
	if err != nil {
		return nil, err
	}
	if num > 0 {
		return nil, cerrors.ErrGroupNameExist
	}
	// 创建group
	groupID := uuid.New()
	if err = query.Group.WithContext(h.Context).Create(&model.Group{
		BaseModel: model.BaseModel{
			ID: groupID,
		},
		GroupName: req.Name,
		OwnerId:   userToken.UserId,
		// TODO: avatar,简介
	}); err != nil {
		hlog.Errorf("err create group: %v", err)
		return nil, err
	}
	// 用户作为owner加入群组
	if err = query.GroupMember.WithContext(h.Context).Create(&model.GroupMember{
		GroupID: groupID,
		UserID:  userToken.UserId,
		Role:    model.GroupRoleAdmin,
	}); err != nil {
		hlog.Errorf("err create group member: %v", err)
		return nil, err
	}
	resp = &group.CreateGroupResp{
		GroupID: groupID.String(),
	}
	return resp, nil
}
