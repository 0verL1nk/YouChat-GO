package cerrors

import "github.com/pkg/errors"

var (
	ErrClientClosed = errors.New("client already closed")
)
