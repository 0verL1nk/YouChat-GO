package cerrors

import "github.com/pkg/errors"

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
