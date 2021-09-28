package user

import (
	"Moreover/pkg/mysql"
	"Moreover/pkg/response"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func Register(id, password string) int {
	if len(password) < 6 || len(id) != 8 {
		return response.ParamError
	}
	if IsUserExist(id) {
		return response.UserExist
	}
	hashPassword, hashError := GenerateHashPassword(password)
	if hashError != nil {
		return response.ERROR
	}
	if err := RegisterUser(id, hashPassword); err != nil {
		return response.ERROR
	}
	return response.SUCCESS
}

func IsUserExist(id string) bool {
	var count int
	sql := `SELECT COUNT(student_id) FROM user WHERE student_id = ?`
	if err := mysql.DB.Get(&count, sql, id); err != nil {
		fmt.Printf("check user exist fail, err: %v\n", err)
		return true
	}
	if count > 0 {
		return true
	}
	return false
}

func GenerateHashPassword(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		fmt.Printf("generate hashPassword fail, err: %v\n", err)
		return "", err
	}
	return string(hashPassword), nil
}

func RegisterUser(id, hashPassword string) error {
	sqlRegister := `INSERT INTO user (student_id, password, permission) VALUES(?, ?, ?)`
	sqlInsertInfo := `INSERT INTO user_info (student_id, nickname, sex, avatar, description) VALUES(?, ?, ?, ?, ?)`
	sql, err := mysql.DB.Begin()
	if err != nil {
		fmt.Printf("register begin fail, err: %v\n", err)
		return err
	}
	if _, err := sql.Exec(sqlRegister, id, hashPassword, 1); err != nil {
		fmt.Printf("insert user fail, err: %v\n", err)
		return err
	}
	if _, err := sql.Exec(sqlInsertInfo, id, DefaultName, DefaultSex, DefaultAvatar, DefaultDescription); err != nil {
		fmt.Printf("insert userInfo fail, err: %v\n", err)
		return err
	}
	if err := sql.Commit(); err != nil {
		fmt.Printf("register commit fail, err: %v\n", err)
		return err
	}
	return nil
}
