package dal

import (
	"core/biz/dal/mysql"
	"core/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
