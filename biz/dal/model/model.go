package model

import (
	"gorm.io/gorm"
)

var AllModels = []interface{}{
	&User{},
	&File{},
}

type User struct {
	UserId    uint64 `gorm:"primaryKey;not null;column:user_id" json:"userId"`
	Name      string `gorm:"type:varchar(255);default:'';column:name" json:"name"`
	Gender    int    `gorm:"type:int;default:0;column:gender" json:"gender"` // 0: 默认，可根据业务调整
	Language  string `gorm:"type:varchar(255);default:'';column:language" json:"language"`
	City      string `gorm:"type:varchar(255);default:'';column:city" json:"city"`
	Province  string `gorm:"type:varchar(255);default:'';column:province" json:"province"`
	Country   string `gorm:"type:varchar(255);default:'';column:country" json:"country"`
	Avatar    string `gorm:"type:varchar(1024);default:'';column:avatar" json:"avatar"`
	Phone     string `gorm:"type:varchar(255);default:'';column:phone" json:"phone"`
	Email     string `gorm:"type:varchar(255);default:'';column:email" json:"email"`
	Password  string `gorm:"type:varchar(255);default:'';column:password" json:"password"`
	Status    int32  `gorm:"type:int;default:0;column:status" json:"status"` // 0: normal, 1: admin, 2: baned
	CreatedAt int64  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt int64  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	IsDeleted bool   `gorm:"column:is_deleted;default:false" json:"is_deleted"`
	IsAdmin   bool   `gorm:"column:is_admin;default:false" json:"is_admin"`
}

type File struct {
	gorm.Model
	UserId         uint64 `gorm:"type:bigint;not null"`
	FileOriginName string `gorm:"type:varchar(255);not null"`
	Key            string `gorm:"type:varchar(255);not null"`
	Usage          string `gorm:"type:varchar(255);not null"`
}

type Querier interface {
	// GetByID(id uint64) (*Event, error)
	// CreateEvent(event *Event) error
	// ListEvents() ([]*Event, error)
	// GetUserInfoByUserId(UserId string) (*User, error)
	// CreateUser(user *User) error
	// ListUsers() ([]*User, error)

	// SELECT * FROM @@table WHERE user_id = @userId AND is_deleted is not true
	GetUserInfoByUserId(userId uint64) (*User, error)

	// SELECT * FROM @@table WHERE email = @email AND is_deleted is not true
	GetUserInfoByEmail(email string) (*User, error)
}
