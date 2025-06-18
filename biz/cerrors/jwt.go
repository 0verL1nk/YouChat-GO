package cerrors

import "github.com/pkg/errors"

var (
	ErrGetTokenFromMiddleware    = errors.New("get token from middleware failed")
	ErrTokenNotFoundInMiddleware = errors.New("token not found in middleware")
	ErrAssertTokenClaims         = errors.New("assert token claims failed")
	ErrParseJWTtoken             = errors.New("parse jwt token failed")
)
