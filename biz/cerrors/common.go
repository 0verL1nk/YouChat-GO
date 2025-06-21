package cerrors

import "github.com/pkg/errors"

var (
	ErrParseUUID             = errors.New("parse uuid failed")
	ErrFetchDataFromDatabase = errors.New("fetch data from database failed")
)
