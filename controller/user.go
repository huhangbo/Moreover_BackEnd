package controller

import (
	"Moreover/internal/pkg/user"
	"Moreover/pkg/jwt"
	"Moreover/pkg/response"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	id := c.PostForm("username")
	password := c.PostForm("password")
	code := user.Register(id, password)
	if code != response.SUCCESS {
		response.Response(c, code, nil)
		return
	}
	token := jwt.GenerateToken(id)
	response.Response(c, code, gin.H{
		"token": token,
	})
}

func Login(c *gin.Context) {
	id := c.PostForm("username")
	password := c.PostForm("password")
	code := user.Login(id, password)
	if code != response.SUCCESS {
		response.Response(c, code, nil)
		return
	}
	token := jwt.GenerateToken(id)
	response.Response(c, code, gin.H{
		"token": token,
	})
}
