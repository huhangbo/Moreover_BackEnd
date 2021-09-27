package model

import "gorm.io/gorm"

type UserInfo struct {
	gorm.Model
	StudentID string
	Nickname string
	Telephone string
	Sex string
	Description string
	Avatar string
}