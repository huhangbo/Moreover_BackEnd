package dao

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

type Comment struct {
	PublishedAt int64                 `gorm:"-" json:"publishedAt"`
	CreatedAt   time.Time             `json:"-"`
	UpdatedAt   time.Time             `json:"-"`
	DeletedAt   soft_delete.DeletedAt `json:"-"`
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
	IsStar        bool          `json:"isStar"`
}

type CommentChild struct {
	CommentDetail
	ReplierInfo UserInfoBasic `json:"replierInfo"`
}
