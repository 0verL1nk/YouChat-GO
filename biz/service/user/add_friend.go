package user

import (
	"context"

	user "core/hertz_gen/user"
	"github.com/cloudwego/hertz/pkg/app"
)

type AddFriendService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewAddFriendService(Context context.Context, RequestContext *app.RequestContext) *AddFriendService {
	return &AddFriendService{RequestContext: RequestContext, Context: Context}
}

func (h *AddFriendService) Run(req *user.AddFriendReq) (resp *user.AddFriendResp, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code
	// 添加好友时创建一个群组
	return
}
