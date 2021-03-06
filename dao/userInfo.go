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
	Background  string                `gorm:"default:http://www.cumt.edu.cn/_upload/article/images/e6/c7/4bd9beeb4438899f538acd1d3360/6a204d99-cff9-41a6-8f09-6dbdd6d34f13.jpg" json:"background"`
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

type UserInfoBasicFollow struct {
	UserInfoBasic
	IsFollow bool `json:"isFollow"`
}
