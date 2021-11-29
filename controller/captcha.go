package controller

import (
	"Moreover/conn"
	"Moreover/pkg/response"
	"Moreover/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
)

var uid uuid.UUID

func GenerateCaptcha(c *gin.Context) {
	id, base64 := service.GenerateCaptcha()
	if id == "" {
		response.Response(c, response.ERROR, nil)
		return
	}
	uid = uuid.New()
	conn.Redis.Set("captcha:"+uid.String(), id, time.Minute*5)
	response.Response(c, response.SUCCESS, gin.H{
		"id":     id,
		"base64": base64,
	})
}

func ParseCaptcha(c *gin.Context) {
	requestId := c.PostForm("captcha")
	id, err := conn.Redis.Get("captcha:" + uid.String()).Result()
	if err != nil {
		response.Response(c, response.ERROR, nil)
		return
	}
	if service.ParseCaptcha(id, requestId) {
		response.Response(c, response.SUCCESS, nil)
		return
	}
	response.Response(c, response.ParamError, nil)
}
