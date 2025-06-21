package user

import (
	"context"

	"core/biz/dal/model"
	"core/biz/dal/query"
	"core/biz/utils"
	user "core/hertz_gen/user"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type GetUserContactsService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetUserContactsService(Context context.Context, RequestContext *app.RequestContext) *GetUserContactsService {
	return &GetUserContactsService{RequestContext: RequestContext, Context: Context}
}

func (h *GetUserContactsService) Run(req *user.GetUserContactsReq) (resp *user.GetUserContactsResp, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code
	// 初始化参数
	// page, pageSize := utils.GetDefaultPageParam(req.Page, req.PageSize)
	// 获取userID
	userToken, err := utils.GetTokenFromMiddleware(h.RequestContext)
	if err != nil {
		return
	}
	// 确认用户状态
	if err = CheckUserExist(h.Context, userToken.UserId); err != nil {
		return
	}
	// 获取用户加入的group
	groups, err := query.Group.GetUserGroups(userToken.UserId)
	if err != nil {
		hlog.Debug("err get user groups:", err)
		return nil, err
	}
	resp = constructConversationResp(groups)
	return resp, nil
}

func constructConversationResp(groups []*model.Group) (resp *user.GetUserContactsResp) {
	var contacts []*user.Contact
	for _, g := range groups {
		contacts = append(contacts, &user.Contact{
			GroupId: g.ID.String(),
			Name:    g.GroupName,
			Avatar:  g.Avatar,
			// 上条消息时间,未读消息数
		})
	}
	return &user.GetUserContactsResp{
		Contacts: contacts,
	}
}
