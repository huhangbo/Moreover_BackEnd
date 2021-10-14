package dao

import (
	"Moreover/connent"
	"gorm.io/plugin/soft_delete"
	"time"
)

type Model interface {
}

type User struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt soft_delete.DeletedAt
	StudentId string `gorm:"primaryKey" json:"student_id"`
	Password  string `json:"password"`
}

func (t User) Add() {
	connent.MySQL.Create(&t)
}

func (t *User) Get() {
	connent.MySQL.First(t, &t.StudentId)
}

func (t User) Delete() {
	connent.MySQL.Delete(&t)
}

func (t User) Update() {
	connent.MySQL.Model(&t).Updates(t)
}
