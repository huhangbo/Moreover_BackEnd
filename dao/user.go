package dao

import (
	"Moreover/conn"
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

func (t User) Add() {
	conn.MySQL.Create(&t)
}

func (t *User) Get() {
	conn.MySQL.First(t, &t.StudentId)
}

func (t User) Delete() {
	conn.MySQL.Delete(&t)
}

func (t User) Update() {
	conn.MySQL.Model(&t).Updates(t)
}
