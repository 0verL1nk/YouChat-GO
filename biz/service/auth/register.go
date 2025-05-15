package auth

import (
	"context"
	"core/biz/dal/model"
	"core/biz/dal/query"
	"core/biz/dal/redis"
	"core/biz/utils"
	auth "core/hertz_gen/auth"
	"errors"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/gorm"
)

var (
	ErrUserNoFound      = errors.New("用户不存在")
	ErrUserAlreadyExist = errors.New("用户已存在")
	ErrUserProhibit     = errors.New("用户被禁用")
	ErrUserCreate       = errors.New("用户创建失败")
	ErrHashPwd          = errors.New("密码加密失败")
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
	if !redis.RedisVerify(h.Context, req.CheckCodeKey, req.CheckCode, true) {
		return &auth.RegisterResp{}, errors.New("验证码错误")
	}
	_, err = CheckUserState(h.Context, req.Email)
	hlog.Debug("err:", err)
	if err == nil || errors.Is(err, ErrUserProhibit) {
		return &auth.RegisterResp{}, ErrUserAlreadyExist
	}
	if !errors.Is(err, ErrUserNoFound) {
		hlog.Error("check user state failed:", err)
		return &auth.RegisterResp{}, err
	}
	pwd, err := utils.HashAndSalt(req.Password)
	if err != nil {
		hlog.Error("hash password failed:", err)
		return &auth.RegisterResp{}, ErrHashPwd
	}
	user := &model.User{
		Email:    req.Email,
		Password: pwd,
		Name:     req.NickName,
	}
	err = query.Q.User.Create(user)
	if err != nil {
		hlog.Error("create user failed:", err)
		return &auth.RegisterResp{}, ErrUserCreate
	}
	return &auth.RegisterResp{Info: "Success"}, nil
}

func CheckUserState(ctx context.Context, email string) (user *model.User, err error) {
	user, err = query.Q.User.GetUserInfoByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &model.User{}, ErrUserNoFound
		}
		return &model.User{}, err
	}
	if user.Status == 2 || user.IsDeleted {
		return &model.User{}, ErrUserProhibit
	}
	return
}
