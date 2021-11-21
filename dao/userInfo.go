package dao

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

type UserInfo struct {
	CreatedAt   time.Time             `json:"-"`
	UpdatedAt   time.Time             `json:"-"`
	DeletedAt   soft_delete.DeletedAt `json:"-"`
	StudentId   string                `gorm:"primaryKey"`
	Nickname    string                `gorm:"default:娶个名字吧"`
	Avatar      string                `gorm:"default:https://moreover-1305054989.cos.ap-nanjing.myqcloud.com/author.jpg"`
	Sex         string                `gorm:"default:未知"`
	Description string                `gorm:"default:添加一句话描述自己吧"`
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
