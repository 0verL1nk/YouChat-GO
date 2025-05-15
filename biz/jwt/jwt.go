package jwt

import (
	"context"
	"core/conf"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/golang-jwt/jwt/v5"
)

var config = conf.GetConf()

var (
	ErrMissingToken = errors.New("missing token")
)

type TokenClaims struct {
	UserId uint64
	jwt.RegisteredClaims
}

func JwtMiddleware() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		authorization := ctx.GetHeader("Authorization")
		authString := string(authorization)
		parts := strings.Split(authString, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			token := parts[1]
			claims, err := ParseToken(token)
			if err != nil {
				ctx.JSON(http.StatusUnauthorized, map[string]string{
					"error": "token is invalid"})
				ctx.Abort()
				return
			}
			ctx.Set("token", claims)
			ctx.Next(c)
		} else {
			ctx.JSON(http.StatusUnauthorized, map[string]string{
				"error": "token is missing"})
			ctx.Abort()
			return
		}
	}

}

// tokenCreate
func CreateToken(ctx context.Context, userId uint64) (token string, err error) {
	tokenClaims := &TokenClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * time.Duration(config.JWT.ValidDays))),
		},
	}
	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims).SignedString([]byte(config.JWT.Secret))
	if err != nil {
		return "", err
	}
	return
}

// parse token

func ParseToken(token string) (tokenClaims *TokenClaims, err error) {
	_, err = jwt.ParseWithClaims(token, tokenClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWT.Secret), nil
	})
	if err != nil {
		return &TokenClaims{}, err
	}
	return
}

// get token from headers
func AuthToken(c context.Context, ctx *app.RequestContext) (user *TokenClaims, err error) {
	// get token from headers
	authHeader := ctx.Request.Header.Get("Authorization")
	if authHeader == "" {
		return &TokenClaims{}, ErrMissingToken
	}
	// token 格式为 "Bearer <token>"
	token := strings.TrimPrefix(authHeader, "Bearer ")
	user, err = ParseToken(token)
	if err != nil {
		return &TokenClaims{}, err
	}
	return
}
