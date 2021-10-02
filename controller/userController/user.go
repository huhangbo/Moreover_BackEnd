package userController

import (
	"Moreover/internal/pkg/user"
	"Moreover/model"
	"Moreover/pkg/jwt"
	"Moreover/pkg/response"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var tmpUser model.User
	if err := c.BindJSON(&tmpUser); err != nil {
		response.Response(c, response.ERROR, nil)
		return
	}
	code := user.Register(tmpUser.UserName, tmpUser.Password)
	if code != response.SUCCESS {
		response.Response(c, code, nil)
		return
	}
	token := jwt.GenerateToken(tmpUser.UserName)
	response.Response(c, code, gin.H{
		"token": token,
	})
}

func Login(c *gin.Context) {
	var tmpUser model.User
	if err := c.BindJSON(&tmpUser); err != nil {
		response.Response(c, response.ERROR, nil)
		return
	}
	code := user.Login(tmpUser.UserName, tmpUser.Password)
	if code != response.SUCCESS {
		response.Response(c, code, nil)
		return
	}
	token := jwt.GenerateToken(tmpUser.UserName)
	response.Response(c, code, gin.H{
		"token": token,
	})
}
