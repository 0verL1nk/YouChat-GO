package auth

import (
	"context"
	"errors"

	"core/biz/dal/redis"
	"core/biz/jwt"
	"core/biz/utils"
	auth "core/hertz_gen/auth"

	"fmt"

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
	if !redis.RedisVerify(h.Context, req.CheckCodeKey, req.CheckCode, true) {
		return &auth.LoginResp{}, errors.New("验证码错误")
	}
	user, err := CheckUserState(h.Context, req.Email)
	if err != nil {
		return &auth.LoginResp{}, err
	}
	err = utils.ComparePasswords(user.Password, req.Password)
	if err != nil {
		return &auth.LoginResp{}, ErrValidatePwd
	}
	token, err := jwt.CreateToken(h.Context, user.UserId)
	if err != nil {
		hlog.Error("create token failed:", err)
		return &auth.LoginResp{}, ErrCreateToken
	}
	return &auth.LoginResp{
		Info:     "Success",
		Token:    token,
		UserId:   fmt.Sprint(user.UserId),
		NickName: user.Name,
		Admin:    user.IsAdmin,
	}, nil
}
