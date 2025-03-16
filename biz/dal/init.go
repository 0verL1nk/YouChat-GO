package dal

import (
	"core/biz/dal/mysql"
)

func Init() {
	// redis.Init()
	mysql.Init()
}
