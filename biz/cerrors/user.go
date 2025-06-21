package cerrors

import "github.com/pkg/errors"

var (
	ErrUserNoFound      = errors.New("user not found")
	ErrUserProhibit     = errors.New("user is prohibited")
	ErrUserAlreadyExist = errors.New("user already exists")
	ErrHashPwd          = errors.New("hash password failed")
	ErrUserCreate       = errors.New("user create failed")
	ErrUserNotOnline    = errors.New("user not online")
)
