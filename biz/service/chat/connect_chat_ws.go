package chat

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

type ConnectChatWSService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}
