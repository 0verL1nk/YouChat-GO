package cerrors

import "github.com/pkg/errors"

var (
	ErrUserNotInGroup = errors.New("user not in group")
)
