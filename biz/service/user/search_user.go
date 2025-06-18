package user

import (
	"context"

	user "core/hertz_gen/user"
	"github.com/cloudwego/hertz/pkg/app"
)

type SearchUserService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewSearchUserService(Context context.Context, RequestContext *app.RequestContext) *SearchUserService {
	return &SearchUserService{RequestContext: RequestContext, Context: Context}
}

func (h *SearchUserService) Run(req *user.SearchUserReq) (resp *user.SearchUserResp, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code
	return
}
