package user

import (
	"Moreover/pkg/mysql"
	"Moreover/pkg/response"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func Login(id, password string) int {
	if !IsUserExist(id) {
		return response.UserNotExist
	}
	var hashPassword string
	sql := `SELECT password FROM user WHERE student_id = ?`
	if err := mysql.DB.Get(&hashPassword, sql, id); err != nil {
		fmt.Printf("get hashPassword fail, err: %v\n", err)
		panic(err)
		return response.ERROR
	}
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	if err != nil {
		fmt.Printf("wrong password\n")
		return response.PasswordError
	}
	return response.SUCCESS
}
