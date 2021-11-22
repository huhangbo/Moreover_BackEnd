package dao

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

type Follow struct {
	ID        uint                  `gorm:"autoIncrement primaryKey" json:"-"`
	CreatedAt time.Time             `json:"createdAt"`
	DeletedAt soft_delete.DeletedAt `json:"-"`
	Parent    string                `json:"parent"`
	Publisher string                `json:"publisher"`
}
