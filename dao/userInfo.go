package dao

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

type UserInfo struct {
	ID          uint                  `gorm:"autoIncrement primaryKey" json:"-"`
	CreatedAt   time.Time             `json:"-"`
	UpdatedAt   time.Time             `json:"-"`
	DeletedAt   soft_delete.DeletedAt `json:"-"`
	StudentId   string                `gorm:"primaryKey" json:"studentId"`
	Nickname    string                `gorm:"default:取个名字吧" json:"nickname"`
	Avatar      string                `gorm:"default:https://moreover-1305054989.cos.ap-nanjing.myqcloud.com/author.jpg" json:"avatar"`
	Sex         string                `gorm:"default:未知" json:"sex"`
	Description string                `gorm:"default:添加一句话描述自己吧" json:"description"`
}

type UserInfoBasic struct {
	StudentId string `gorm:"primaryKey" json:"studentId"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	Sex       string `json:"sex"`
}

type UserInfoDetail struct {
	UserInfo
	Follower int  `json:"follower"`
	Fan      int  `json:"fan"`
	IsFollow bool `json:"isFollow"`
}
