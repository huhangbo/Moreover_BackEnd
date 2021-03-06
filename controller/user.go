package controller

import (
	"Moreover/dao"
	"Moreover/pkg/jwt"
	"Moreover/pkg/response"
	"Moreover/service"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var tmpUser dao.User
	if err := c.BindJSON(&tmpUser); err != nil {
		response.Response(c, response.ERROR, nil)
		return
	}
	code := service.Register(tmpUser)
	if code != response.SUCCESS {
		response.Response(c, code, nil)
		return
	}
	token := jwt.GenerateToken(tmpUser.StudentId)
	response.Response(c, code, gin.H{
		"token": token,
	})
}

func Login(c *gin.Context) {
	var tmpUser dao.User
	if err := c.BindJSON(&tmpUser); err != nil {
		response.Response(c, response.ERROR, nil)
		return
	}
	code := service.Login(tmpUser)
	if code != response.SUCCESS {
		response.Response(c, code, nil)
		return
	}
	token := jwt.GenerateToken(tmpUser.StudentId)
	response.Response(c, code, gin.H{
		"token": token,
	})
}
