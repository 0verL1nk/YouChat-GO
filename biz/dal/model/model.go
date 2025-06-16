package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var AllModels = []interface{}{
	&User{},
	&File{},
	&ChatMessage{},
	&Group{},
}

type User struct {
	gorm.Model
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

type Group struct {
	gorm.Model
	GroupName string  `gorm:"type:varchar(255);default:'';column:group_name" json:"groupName"`
	OwnerId   uint    `gorm:"type:bigint;not null;column:owner_id" json:"ownerId"`
	Avatar    string  `gorm:"type:varchar(1024);default:'';column:avatar" json:"avatar"`
	Member    []*User `gorm:"many2many:group_member;"`
	Desc      string  `gorm:"type:varchar(1024);default:'';column:desc" json:"desc"` //群组简介
}

type Conversation struct {
	gorm.Model
	UserID      uint `gorm:"type:bigint;column:user_id"`
	GroupID     uint `gorm:"type:bigint;column:group_id"`
	LastMsgTime time.Time
	UnReadNum   uint
}

type ChatMessage struct {
	gorm.Model
	MsgType MessageType `gorm:"type:int;not null;default:0;column:msg_type" json:"msg_type"`
	Content string      `gorm:"type:varchar(10240);default:'';column:content" json:"content"`
	FromId  uint        `gorm:"type:bigint;not null"`
	ToId    uint        `gorm:"type:bigint;not null"`
}

type MessageType uint8

const (
	MsgTypeText MessageType = iota
	MsgTypeImage
	MsgTypeAudio
	MsgTypeVideo
	MsgTypeFile
)

func (m MessageType) String() string {
	switch m {
	case MsgTypeText:
		return "text"
	case MsgTypeImage:
		return "image"
	case MsgTypeAudio:
		return "audio"
	case MsgTypeVideo:
		return "video"
	case MsgTypeFile:
		return "file"
	default:
		return "unknown"
	}
}

var Str2MsgType = map[string]MessageType{
	"text":  MsgTypeText,
	"image": MsgTypeImage,
	"audio": MsgTypeAudio,
	"video": MsgTypeVideo,
	"file":  MsgTypeFile,
}

type File struct {
	gorm.Model
	UserId         uint64    `gorm:"type:bigint;not null"`
	SID            uuid.UUID `gorm:"column:sid"`
	FileOriginName string    `gorm:"type:varchar(255);not null"`
	Key            string    `gorm:"type:varchar(255);not null"`
	Usage          string    `gorm:"type:varchar(255);not null"`
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
