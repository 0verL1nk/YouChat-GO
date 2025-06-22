package cerrors

import "github.com/pkg/errors"

var (
	BaseErrBadReq   = errors.New("bad request")
	BaseErrInternal = errors.New("internal server error")
	BaseErrNotFound = errors.New("not found")
)

func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	return errors.Wrap(err, msg)
}

func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return errors.Wrapf(err, format, args...)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}
