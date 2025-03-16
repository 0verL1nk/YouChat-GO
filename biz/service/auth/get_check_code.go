package auth

import (
	"context"
	"core/biz/dal/redis"
	auth "core/hertz_gen/auth"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/mojocn/base64Captcha"
)

var rsStore *redis.RedisStore

type GetCheckCodeService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetCheckCodeService(Context context.Context, RequestContext *app.RequestContext) *GetCheckCodeService {
	return &GetCheckCodeService{RequestContext: RequestContext, Context: Context}
}

func (h *GetCheckCodeService) Run(req *auth.GetCheckCodeReq) (resp *auth.GetCheckCodeResp, err error) {
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
	rsStore.Set(id, ans)
	var data auth.CheckCode
	data.CheckCode = b64s
	data.CheckCodeKey = id
	resp = &auth.GetCheckCodeResp{}
	resp.Data = &data
	resp.Code = 200
	resp.Info = "Success"
	return
}
