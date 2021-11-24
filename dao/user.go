package dao

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

type Model interface {
}

type User struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt soft_delete.DeletedAt
	StudentId string `gorm:"primaryKey" json:"studentId"`
	Password  string `json:"password"`
}
