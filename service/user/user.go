package user

import (
	"Moreover/connent"
	"Moreover/dao"
	"Moreover/pkg/response"
	"golang.org/x/crypto/bcrypt"
)

func Register(user dao.User) int {
	if len(user.Password) < 6 || len(user.StudentId) != 8 {
		return response.ParamError
	}
	if isUserExist(user) {
		return response.UserExist
	}
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashPassword)
	var tmpUserInfo dao.UserInfo
	tmpUserInfo.StudentId = user.StudentId
	tx := connent.MySQL.Begin()
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
	if !isUserExist(user) {
		return response.UserNotExist
	}
	password := user.Password
	connent.MySQL.First(&user, user.StudentId)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return response.PasswordError
	}
	return response.SUCCESS
}

func isUserExist(user dao.User) bool {
	if err := connent.MySQL.First(&user, user.StudentId).Error; err != nil {
		return false
	}
	return true
}
