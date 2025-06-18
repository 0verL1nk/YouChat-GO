package utils

import (
	"strconv"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// 从 c.Query 中读取 T 类型的 param。key: 查询参数；
// def: 如果查询参数不存在或转换失败，返回的默认值。
// 可接受的类型定义在 utils.ParsableQueryParam 接口中
type ParsableQueryParam interface {
	~int | ~int64 | ~float64 | ~string | uuid.UUID | *timestamppb.Timestamp
}

// 从 c.Query 中读取 T 类型的 param。key: 查询参数；
// def: 如果查询参数不存在或转换失败，返回的默认值。
// 可接受的类型定义在 utils.ParsableQueryParam 接口中
func QueryParam[T ParsableQueryParam](c *app.RequestContext, key string, def T) T {
	queryValue := c.Query(key)
	if queryValue == "" {
		return def
	}

	var zero T
	switch any(zero).(type) {
	case int:
		if v, err := strconv.Atoi(queryValue); err == nil {
			return any(v).(T)
		}
	case int64:
		if v, err := strconv.ParseInt(queryValue, 10, 64); err == nil {
			return any(v).(T)
		}
	case float64:
		if v, err := strconv.ParseFloat(queryValue, 64); err == nil {
			return any(v).(T)
		}
	case string:
		return any(queryValue).(T)
	case uuid.UUID:
		if v, err := uuid.Parse(queryValue); err == nil {
			return any(v).(T)
		}
	}

	return def
}
func QueryParamInt(c *app.RequestContext, key string, def int) int {
	return QueryParam(c, key, def)
}

func QueryParamInt64(c *app.RequestContext, key string, def int64) int64 {
	return QueryParam(c, key, def)
}

// 返回正整数
func QueryParamInt64Pos(c *app.RequestContext, key string, def int64) int64 {
	res := QueryParam(c, key, def)
	if res > 0 {
		return res
	} else {
		return def
	}
}
func Atoi64(s string) int64 {
	if v, err := strconv.ParseInt(s, 10, 64); err == nil {
		return v
	}
	return 0
}

func ParseTimestamp(t *timestamppb.Timestamp, def time.Time) time.Time {
	if t == nil {
		return def
	}
	return time.Unix(t.GetSeconds(), 0).In(time.Local)
}

func GetDefaultPageParam(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	return page, pageSize
}
