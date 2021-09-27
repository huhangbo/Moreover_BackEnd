package controller

import (
	"Moreover/internal/pkg/jwt"
	"Moreover/internal/pkg/response"
	"Moreover/internal/pkg/user"
	"github.com/gin-gonic/gin"
)


func Register(c *gin.Context)  {
	id := c.PostForm("username")
	password := c.PostForm("password")
	code := user.Register(id, password)
	if code != response.SUCCESS{
		response.Response(c, code, nil)
		return
	}
	token :=jwt.GetToken(id)
	response.Response(c, code, gin.H{
		"token": token,
	})
}

func Login(c *gin.Context)  {
	id := c.PostForm("username")
	password := c.PostForm("password")
	code := user.Login(id, password)
	if code != response.SUCCESS {
		response.Response(c, code, nil)
		return
	}
	token := jwt.GetToken(id)
	response.Response(c, code, gin.H{
		"token": token,
	})
}
