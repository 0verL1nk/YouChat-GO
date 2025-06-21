package user

import (
	"context"
	"core/biz/cerrors"
	"core/biz/dal/model"
	"core/biz/dal/query"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CheckUserStateByEmail(ctx context.Context, email string) (user *model.User, err error) {
	user, err = query.Q.User.GetUserInfoByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &model.User{}, cerrors.ErrUserNoFound
		}
		return &model.User{}, err
	}
	if user.Status == 2 || user.DeletedAt.Valid {
		return &model.User{}, cerrors.ErrUserProhibit
	}
	return
}

func CheckUserExist(ctx context.Context, userID uuid.UUID) (err error) {
	user, err := query.Q.User.WithContext(ctx).Where(query.User.ID.Eq(userID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return cerrors.ErrUserNoFound
		}
		return err
	}
	if user.Status == 2 || user.DeletedAt.Valid {
		return cerrors.ErrUserProhibit
	}
	return nil
}

// func GetUserGroups(ctx context.Context, userID uuid.UUID, page int, page_size int) (resp []*model.Group, err error) {
// 	if _, err = query.Group.WithContext(ctx).
// 		Join(query.GroupMember, query.GroupMember.UserID.Eq(userID)).
// 		Where(query.Group.ID.Eq(q))).
// 		ScanByPage(&resp); err != nil {
// 		return
// 	}
// 	return
// }
