package utils

import (
	"core/biz/cerrors"
	"core/biz/jwt"

	"github.com/bwmarrin/snowflake"
	"github.com/cloudwego/hertz/pkg/app"
)

// 雪花算法生成uint64的随机ID
func GenNumId() (res int64, err error) {
	node, err := snowflake.NewNode(1)
	if err != nil {
		return 0, err
	}
	return node.Generate().Int64(), nil
}

// jwt中间件后,从请求上下文获取token

func GetTokenFromMiddleware(c *app.RequestContext) (*jwt.TokenClaims, error) {
	t, exists := c.Get("token")
	if !exists {
		return nil, cerrors.ErrTokenNotFoundInMiddleware
	}
	token, ok := t.(*jwt.TokenClaims)
	if !ok {
		return nil, cerrors.ErrAssertTokenClaims
	}
	return token, nil
}
