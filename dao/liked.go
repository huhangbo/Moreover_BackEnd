package dao

import (
	"time"
)

type Liked struct {
	CreatedAt time.Time `json:"createdAt"`
	ParentId  string    `gorm:"primaryKey" json:"parentId"`
	Publisher string    `gorm:"primaryKey" json:"publisher"`
	LikeUser  string    `json:"likeUser"`
}
