package auth

import (
	"context"
	"errors"

	"core/biz/dal/redis"
	"core/biz/jwt"
	"core/biz/service/user"
	"core/biz/utils"
	auth "core/hertz_gen/auth"
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

var (
	ErrCreateToken = errors.New("create token failed")
	ErrValidatePwd = errors.New("账号或密码错误")
)

type LoginService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewLoginService(Context context.Context, RequestContext *app.RequestContext) *LoginService {
	return &LoginService{RequestContext: RequestContext, Context: Context}
}

func (h *LoginService) Run(req *auth.LoginReq) (resp *auth.LoginResp, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code
	// 验证码校验
	if config.Hertz.EnableCaptcha {
		if !redis.RedisVerify(h.Context, req.CaptchaKey, req.Captcha, true) {
			return &auth.LoginResp{}, errors.New("验证码错误")
		}
	}
	user, err := user.CheckUserStateByEmail(h.Context, req.Email)
	if err != nil {
		return &auth.LoginResp{}, err
	}
	err = utils.ComparePasswords(user.Password, req.Password)
	if err != nil {
		return &auth.LoginResp{}, ErrValidatePwd
	}
	token, expireAt, err := jwt.CreateToken(h.Context, user.ID)
	if err != nil {
		hlog.Error("create token failed:", err)
		return &auth.LoginResp{}, ErrCreateToken
	}
	return &auth.LoginResp{
		Token:         token,
		UserId:        fmt.Sprint(user.ID),
		NickName:      user.Name,
		TokenExpireAt: timestamppb.New(expireAt),
	}, nil
}
