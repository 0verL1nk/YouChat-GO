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
	&Conversation{},
	&GroupMember{},
}

type BaseModel struct {
	ID        uuid.UUID `gorm:"type:char(36);primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// BeforeCreate 钩子函数：在创建记录前生成 UUID
func (b *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	if b.CreatedAt.IsZero() {
		b.CreatedAt = time.Now()
	}
	if b.UpdatedAt.IsZero() {
		b.UpdatedAt = time.Now()
	}
	return nil
}

type User struct {
	BaseModel
	Name     string `gorm:"type:varchar(255);default:'';column:name" json:"name"`
	Gender   int    `gorm:"type:int;default:0;column:gender" json:"gender"` // 0: 默认，可根据业务调整
	Language string `gorm:"type:varchar(255);default:'';column:language" json:"language"`
	City     string `gorm:"type:varchar(255);default:'';column:city" json:"city"`
	Province string `gorm:"type:varchar(255);default:'';column:province" json:"province"`
	Country  string `gorm:"type:varchar(255);default:'';column:country" json:"country"`
	Avatar   string `gorm:"type:varchar(1024);default:'';column:avatar" json:"avatar"`
	Phone    string `gorm:"type:varchar(255);default:'';column:phone" json:"phone"`
	Email    string `gorm:"type:varchar(255);default:'';column:email" json:"email"`
	Password string `gorm:"type:varchar(255);default:'';column:password" json:"password"`
	Status   int32  `gorm:"type:int;default:0;column:status" json:"status"` // 0: normal, 1: admin, 2: baned
	IsAdmin  bool   `gorm:"column:is_admin;default:false" json:"is_admin"`
}

type Group struct {
	BaseModel
	GroupName string      `gorm:"type:varchar(255);default:'';column:group_name" json:"groupName"`
	OwnerId   uuid.UUID   `gorm:"type:char(36);not null;column:owner_id" json:"ownerId"`
	Avatar    string      `gorm:"type:varchar(1024);default:'';column:avatar" json:"avatar"`
	Desc      string      `gorm:"type:varchar(1024);default:'';column:desc" json:"desc"`   //群组简介
	Status    GroupStatus `gorm:"type:int;not null;default:0;column:status" json:"status"` // 群组状态 0: normal, 1: closed, 2: deleted, 3: banned
}

type GroupType uint8

const (
	GroupTypePublic  GroupType = iota // 公开群组，任何人都可以加入
	GroupTypePrivate                  // 私有群组，只有邀请的用户可以加入
	GroupTypeSecret                   // 密码群组，用户需要输入密码才能加入
	GroupTypeFriend                   //两人群组,好友之间
)

func (g GroupType) String() string {
	switch g {
	case GroupTypePublic:
		return "public"
	case GroupTypePrivate:
		return "private"
	case GroupTypeSecret:
		return "secret"
	case GroupTypeFriend:
		return "friend"
	default:
		return "unknown"
	}
}

var Str2GroupType = map[string]GroupType{
	"public":  GroupTypePublic,
	"private": GroupTypePrivate,
	"secret":  GroupTypeSecret,
	"friend":  GroupTypeFriend,
}

type GroupRole uint8

const (
	GroupRoleOwner  GroupRole = iota // 群主
	GroupRoleAdmin                   // 管理员
	GroupRoleMember                  // 普通成员
)

func (g GroupRole) String() string {
	switch g {
	case GroupRoleOwner:
		return "owner"
	case GroupRoleAdmin:
		return "admin"
	case GroupRoleMember:
		return "member"
	default:
		return "unknown"
	}
}

var Str2GroupRole = map[string]GroupRole{
	"owner":  GroupRoleOwner,
	"admin":  GroupRoleAdmin,
	"member": GroupRoleMember,
}

type GroupStatus uint8

const (
	GroupStatusNormal  GroupStatus = iota // 正常状态
	GroupStatusClosed                     // 已关闭
	GroupStatusDeleted                    // 已删除
	GroupStatusBanned                     // 已封禁
)

func (g GroupStatus) String() string {
	switch g {
	case GroupStatusNormal:
		return "normal"
	case GroupStatusClosed:
		return "closed"
	case GroupStatusDeleted:
		return "deleted"
	case GroupStatusBanned:
		return "banned"
	default:
		return "unknown"
	}
}

var Str2GroupStatus = map[string]GroupStatus{
	"normal":  GroupStatusNormal,
	"closed":  GroupStatusClosed,
	"deleted": GroupStatusDeleted,
	"banned":  GroupStatusBanned,
}

type GroupMember struct {
	BaseModel
	GroupType GroupType `gorm:"type:int;not null;default:0;column:group_type" json:"groupType"` // 群组类型
	GroupID   uuid.UUID `gorm:"type:char(36);not null;column:group_id" json:"groupId"`
	UserID    uuid.UUID `gorm:"type:char(36);not null;column:user_id" json:"userId"`
	Group     Group     `gorm:"foreignKey:GroupID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"group"` // 群组信息
	Role      GroupRole `gorm:"type:int;not null;default:0;column:role" json:"role"`                                        // 群组角色 0: 群主, 1: 管理员, 2: 成员
	Status    uint8     `gorm:"type:int;default:0;column:status" json:"status"`                                             // 0: normal, 1: admin, 2: baned
}

type Conversation struct {
	BaseModel
	UserID      uuid.UUID `gorm:"type:char(36);column:user_id"`
	GroupID     uuid.UUID `gorm:"type:char(36);column:group_id"`
	LastMsgTime time.Time
	UnReadNum   uint64
}

type ChatMessage struct {
	BaseModel
	MsgType MessageType `gorm:"type:int;not null;default:0;column:msg_type" json:"msg_type"`
	Content string      `gorm:"type:varchar(10240);default:'';column:content" json:"content"`
	FromId  uuid.UUID   `gorm:"type:char(36);not null"`
	// 找群组消息直接查ToID
	ToId uuid.UUID `gorm:"type:char(36);not null"`
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
	BaseModel
	UserId         uuid.UUID `gorm:"type:char(36);not null"`
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

	// SELECT * FROM @@table WHERE user_id = @userId
	GetUserInfoByUserId(userId uuid.UUID) (*User, error)

	// SELECT * FROM @@table WHERE email = @email
	GetUserInfoByEmail(email string) (*User, error)

	// SELECT `groups`.* FROM `groups` JOIN `group_members` ON  `group_members`.`user_id`=@userId WHERE `groups`.`id` = `group_members`.`group_id`
	GetUserGroups(userId uuid.UUID) ([]*Group, error)
}
