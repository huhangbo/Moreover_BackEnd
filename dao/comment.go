package dao

import (
	"Moreover/model"
	"gorm.io/plugin/soft_delete"
	"time"
)

type Comment struct {
	CreatedAt time.Time             `json:"createdAt"`
	UpdatedAt time.Time             `json:"-"`
	DeletedAt soft_delete.DeletedAt `json:"-"`
	ParentId  string                `json:"parentId"`
	CommentId string                `gorm:"primaryKey" json:"CommentId"`
	Publisher string                `json:"publisher"`
	Replier   string                `json:"replier"`
	Message   string                `json:"message"`
}

type CommentDetail struct {
	Comment
	PublisherInfo UserInfoBasic
	Star          int  `json:"star"`
	IsStart       bool `json:"isStart"`
}

type CommentChild struct {
	CommentDetail
	ReplierInfo UserInfoBasic
}

type ParentComment struct {
	CommentDetail
	Page model.Page
}
