package group

import (
	"context"
	"errors"

	"core/biz/cerrors"
	"core/biz/dal/model"
	"core/biz/dal/query"
	"core/biz/service/user"
	"core/biz/utils"
	group "core/hertz_gen/group"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JoinGroupService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewJoinGroupService(Context context.Context, RequestContext *app.RequestContext) *JoinGroupService {
	return &JoinGroupService{RequestContext: RequestContext, Context: Context}
}

func (h *JoinGroupService) Run(req *group.JoinGroupReq) (resp *group.JoinGroupResp, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code
	// parse uuid
	_groupID, err := uuid.Parse(req.GroupId)
	if err != nil {
		return nil, cerrors.ErrParseUUID
	}
	// 校验用户
	tokenClaim, err := utils.GetTokenFromMiddleware(h.RequestContext)
	if err != nil {
		return nil, cerrors.ErrGetTokenFromMiddleware
	}
	if err = user.CheckUserExist(h.Context, tokenClaim.UserId); err != nil {
		return nil, err
	}
	// 校验所加入的群组是否合法
	if err = CheckJoinGroup(h.Context, _groupID, tokenClaim.UserId); err != nil {
		return nil, err
	}
	// 加入群组
	if err = query.GroupMember.WithContext(h.Context).Create(&model.GroupMember{
		GroupID: _groupID,
		UserID:  tokenClaim.UserId,
	}); err != nil {
		return nil, err
	}
	return &group.JoinGroupResp{
		GroupId: _groupID.String(),
	}, nil
}

// 校验所加入群组是否合法,群组是否存在,用户是否已经加入
func CheckJoinGroup(ctx context.Context, groupID uuid.UUID, userID uuid.UUID) (err error) {
	group, err := query.Group.WithContext(ctx).Where(query.Group.ID.Eq(groupID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return cerrors.ErrGroupNotFound
		}
		return err
	}
	if group.Status != model.GroupStatusNormal || group.DeletedAt.Valid {
		return cerrors.ErrGroupProhibit
	}
	// 用户是否已经加入
	count, err := query.GroupMember.WithContext(ctx).Where(query.GroupMember.UserID.Eq(userID), query.GroupMember.GroupID.Eq(groupID)).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return cerrors.ErrGroupAlreadyMember
	}
	return nil
}
