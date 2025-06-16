package auth

import (
	"context"

	"core/biz/dal/redis"
	auth "core/hertz_gen/auth"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/mojocn/base64Captcha"
)

type GetCaptchaService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetCaptchaService(Context context.Context, RequestContext *app.RequestContext) *GetCaptchaService {
	return &GetCaptchaService{RequestContext: RequestContext, Context: Context}
}

func (h *GetCaptchaService) Run(req *auth.GetCaptchaReq) (resp *auth.GetCaptchaResp, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code
	c := base64Captcha.NewCaptcha(base64Captcha.NewDriverDigit(42, 95, 4, 0.5, 40), base64Captcha.DefaultMemStore)
	id, b64s, ans, err := c.Generate()
	if err != nil {
		hlog.Error("generate checkCode failed:", err)
		return
	}
	// Set data
	if err = redis.RedisSet(h.Context, id, ans); err != nil {
		hlog.Error("set checkCode failed:", err)
		return
	}
	resp = &auth.GetCaptchaResp{
		Captcha:    b64s,
		CaptchaKey: id,
	}

	return
}
