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
	"github.com/google/uuid"
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
	// 校验用户信息
	if err = userValidator(h.Context, req.Email, req.Password, req.NickName); err != nil {
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
	// 加入公共群聊
	if err = JoinPublicGroup(h.Context, user.ID); err != nil {
		hlog.Error("join public group failed:", err)
	}
	return &auth.RegisterResp{Info: "Success"}, nil
}

func JoinPublicGroup(ctx context.Context, userID uuid.UUID) (err error) {
	groupID, err := query.Group.WithContext(ctx).Where(query.Group.GroupName.Eq("public")).First()
	if err != nil {
		hlog.Error("get public group failed:", err)
		return err
	}
	if _, err = query.GroupMember.WithContext(ctx).
		Where(query.GroupMember.UserID.Eq(userID),
			query.GroupMember.GroupID.Eq(groupID.ID),
			query.GroupMember.GroupType.Eq(uint8(model.GroupTypePublic))).
		FirstOrCreate(); err != nil {
		hlog.Error("join public group failed:", err)
		return err
	}
	return nil
}

func userValidator(ctx context.Context, email, password, nickName string) (err error) {
	if email == "" {
		return cerrors.ErrEmailEmpty
	}
	if password == "" {
		return cerrors.ErrPasswordEmpty
	}
	if !utils.IsValidEmail(email) {
		return cerrors.ErrEmailFormat
	}
	if len(password) < 6 || len(password) > 20 {
		return cerrors.ErrPasswordLength
	}
	if nickName == "" {
		return cerrors.ErrNickNameEmpty
	}
	// 确认昵称是否存在
	num, err := query.Q.User.WithContext(ctx).Where(query.User.Name.Eq(nickName)).Count()
	if err != nil {
		hlog.Error("check nickname exist failed:", err)
		return err
	}
	if num > 0 {
		return cerrors.ErrUserAlreadyExist
	}
	return nil
}
