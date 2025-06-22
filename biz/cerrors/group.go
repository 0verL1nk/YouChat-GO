package cerrors

import "github.com/pkg/errors"

var (
	ErrGroupNotFound      = errors.New("group not found")
	ErrGroupProhibit      = errors.New("group prohibit operation")
	ErrGroupAlreadyMember = errors.New("user already in group")
)
