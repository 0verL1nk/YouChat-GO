package chat

import (
	"context"

	chat "core/hertz_gen/chat"
	"github.com/cloudwego/hertz/pkg/app"
)

type GetConversationService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetConversationService(Context context.Context, RequestContext *app.RequestContext) *GetConversationService {
	return &GetConversationService{RequestContext: RequestContext, Context: Context}
}

func (h *GetConversationService) Run(req *chat.GetConversationReq) (resp *chat.GetConversationResponse, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code

	return
}
