package mysql

import (
	"context"
	"core/biz/dal/model"
	"core/biz/dal/query"
	"core/conf"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
	c   = conf.GetConf()
)

func Init() {
	//dsn := "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&tls=%s"

	DB, err = gorm.Open(mysql.Open(fmt.Sprintf(dsn, c.MySQL.Username, c.MySQL.Password, c.MySQL.Host, c.MySQL.Port, c.MySQL.Database, c.MySQL.TLS)),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		},
	)
	if err != nil {
		panic(err)
	}

}

func CreatePublicGroup(ctx context.Context) {
	num, err := query.Group.WithContext(ctx).Where(query.Group.GroupName.Eq("public")).Count()
	if err != nil {
		panic(err)
	}
	if num > 0 {
		return
	}
	if err = query.Group.WithContext(ctx).Create(&model.Group{
		GroupName: "public",
		Desc:      "公共群组",
	}); err != nil {
		panic(err)
	}
}
