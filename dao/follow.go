package dao

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

type Follow struct {
	CreatedAt time.Time             `json:"createdAt"`
	UpdatedAt time.Time             `json:"updatedAt"`
	DeletedAt soft_delete.DeletedAt `json:"deletedAt"`
	Parent    string                `gorm:"primaryKey" json:"parent"`
	Publisher string                `gorm:"primaryKey" json:"publisher"`
}
