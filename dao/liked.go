package dao

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

type Liked struct {
	ID        uint                  `gorm:"autoIncrement primaryKey" json:"-"`
	CreatedAt time.Time             `json:"-"`
	DeletedAt soft_delete.DeletedAt `json:"-"`
	Parent    string                `json:"parent"`
	Publisher string                `json:"publisher"`
	Liker     string                `json:"liker"`
}
