package dao

import (
	"time"
)

type Follow struct {
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Parent    string    `gorm:"primaryKey" json:"parent"`
	Publisher string    `gorm:"primaryKey" json:"publisher"`
}
