package user

import (
	"context"

	user "core/hertz_gen/user"
	"github.com/cloudwego/hertz/pkg/app"
)

type GetUserConversationsService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetUserConversationsService(Context context.Context, RequestContext *app.RequestContext) *GetUserConversationsService {
	return &GetUserConversationsService{RequestContext: RequestContext, Context: Context}
}

func (h *GetUserConversationsService) Run(req *user.GetUserConversationsReq) (resp *user.GetUserConversationsResp, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code
	return
}
