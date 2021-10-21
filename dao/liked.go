package dao

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

type Liked struct {
	CreatedAt time.Time             `json:"createdAt"`
	DeletedAt soft_delete.DeletedAt `json:"deletedAt"`
	ParentId  string                `gorm:"primaryKey" json:"parentId"`
	Publisher string                `gorm:"primaryKey" json:"publisher"`
	LikeUser  string                `json:"likeUser"`
}
