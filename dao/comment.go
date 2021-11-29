package dao

import (
	"Moreover/model"
	"gorm.io/plugin/soft_delete"
	"time"
)

type Comment struct {
	PublishedAt int64                 `gorm:"-" json:"publishedAt"`
	CreatedAt   time.Time             `json:"-"`
	UpdatedAt   time.Time             `json:"-"`
	DeletedAt   soft_delete.DeletedAt `json:"-"`
	Kind        string                `json:"kind"`
	KindId      string                `json:"kindId"`
	ParentId    string                `json:"parentId"`
	CommentId   string                `gorm:"primaryKey" json:"CommentId"`
	Publisher   string                `json:"publisher"`
	Replier     string                `json:"replier"`
	Message     string                `json:"message"`
}

type CommentDetail struct {
	Comment
	PublisherInfo UserInfoBasic `json:"publisherInfo"`
	Star          int           `json:"star"`
	IsStart       bool          `json:"isStart"`
}

type CommentChild struct {
	CommentDetail
	ReplierInfo UserInfoBasic `json:"replierInfo"`
}

type ParentComment struct {
	CommentDetail
	Page model.Page `json:"page"`
}
