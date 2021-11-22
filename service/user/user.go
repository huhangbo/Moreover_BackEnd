package user

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/pkg/response"
	"golang.org/x/crypto/bcrypt"
)

func Register(user dao.User) int {
	if len(user.Password) < 6 || len(user.StudentId) != 8 {
		return response.ParamError
	}
	if IsUserExist(user) {
		return response.UserExist
	}
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashPassword)
	var tmpUserInfo dao.UserInfo
	tmpUserInfo.StudentId = user.StudentId
	tx := conn.MySQL.Begin()
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return response.ERROR
	}
	if err := tx.Create(&tmpUserInfo).Error; err != nil {
		tx.Rollback()
		return response.ERROR
	}
	tx.Commit()
	return response.SUCCESS
}

func Login(user dao.User) int {
	if !IsUserExist(user) {
		return response.UserNotExist
	}
	password := user.Password
	conn.MySQL.First(&user, user.StudentId)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return response.PasswordError
	}
	return response.SUCCESS
}

func IsUserExist(user dao.User) bool {
	if err := conn.MySQL.First(&user, user.StudentId).Error; err != nil {
		return false
	}
	return true
}
