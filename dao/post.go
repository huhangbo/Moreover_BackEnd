package dao

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

type Post struct {
	PublishedAt int64                 `gorm:"-" json:"publishedAt"`
	CreatedAt   time.Time             `json:"createdAt"`
	UpdatedAt   time.Time             `json:"-"`
	DeletedAt   soft_delete.DeletedAt `json:"-"`
	PostId      string                `gorm:"primaryKey" json:"postId"`
	Publisher   string                `json:"publisher"`
	Pictures    []string              `gorm:"-" json:"picture"`
	Picture     string                `json:"-"`
	Detail      string                `json:"detail"`
}

type PostDetail struct {
	Post          `gorm:"embedded"`
	Star          int           `json:"star"`
	IsStar        bool          `json:"isStar"`
	Comments      int           `json:"comments"`
	PublisherInfo UserInfoBasic `json:"publisherInfo"`
}
