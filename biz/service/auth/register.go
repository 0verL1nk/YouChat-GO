package auth

import (
	"context"
	"core/biz/cerrors"
	"core/biz/dal/model"
	"core/biz/dal/query"
	"core/biz/dal/redis"
	"core/biz/service/user"
	"core/biz/utils"
	"core/conf"
	auth "core/hertz_gen/auth"
	"errors"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

var (
	config = conf.GetConf()
)

type RegisterService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewRegisterService(Context context.Context, RequestContext *app.RequestContext) *RegisterService {
	return &RegisterService{RequestContext: RequestContext, Context: Context}
}

func (h *RegisterService) Run(req *auth.RegisterReq) (resp *auth.RegisterResp, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code
	if config.Hertz.EnableCaptcha {
		if !redis.RedisVerify(h.Context, req.CaptchaKey, req.Captcha, true) {
			return &auth.RegisterResp{}, errors.New("验证码错误")
		}
	}
	_, err = user.CheckUserStateByEmail(h.Context, req.Email)
	hlog.Debug("err:", err)
	if err == nil || errors.Is(err, cerrors.ErrUserProhibit) {
		return &auth.RegisterResp{}, cerrors.ErrUserAlreadyExist
	}
	if !errors.Is(err, cerrors.ErrUserNoFound) {
		hlog.Error("check user state failed:", err)
		return &auth.RegisterResp{}, err
	}
	pwd, err := utils.HashAndSalt(req.Password)
	if err != nil {
		hlog.Error("hash password failed:", err)
		return &auth.RegisterResp{}, cerrors.ErrHashPwd
	}
	user := &model.User{
		Email:    req.Email,
		Password: pwd,
		Name:     req.NickName,
	}
	err = query.Q.User.Create(user)
	if err != nil {
		hlog.Error("create user failed:", err)
		return &auth.RegisterResp{}, cerrors.ErrUserCreate
	}
	return &auth.RegisterResp{Info: "Success"}, nil
}
