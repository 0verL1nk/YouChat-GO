package cerrors

import "github.com/pkg/errors"

var (
	ErrUserNoFound      = errors.New("user not found")
	ErrUserProhibit     = errors.New("user is prohibited")
	ErrUserAlreadyExist = errors.Wrap(BaseErrBadReq, "user already exists")
	ErrHashPwd          = errors.New("hash password failed")
	ErrUserCreate       = errors.New("user create failed")
	ErrUserNotOnline    = errors.New("user not online")
	ErrEmailEmpty       = errors.Wrap(BaseErrBadReq, "email cannot be empty")
	ErrPasswordEmpty    = errors.Wrap(BaseErrBadReq, "password cannot be empty")
	ErrEmailFormat      = errors.Wrap(BaseErrBadReq, "email format is invalid")
	ErrPasswordLength   = errors.Wrap(BaseErrBadReq, "password length must be between 6 and 20 characters")
	ErrNickNameEmpty    = errors.Wrap(BaseErrBadReq, "nickname cannot be empty")
)
