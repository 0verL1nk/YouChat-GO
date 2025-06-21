package group

import (
	"context"

	group "core/hertz_gen/group"
	"github.com/cloudwego/hertz/pkg/app"
)

type SearchGroupService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewSearchGroupService(Context context.Context, RequestContext *app.RequestContext) *SearchGroupService {
	return &SearchGroupService{RequestContext: RequestContext, Context: Context}
}

func (h *SearchGroupService) Run(req *group.SearchGroupReq) (resp *group.SearchGroupResp, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code
	return
}
