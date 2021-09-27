package controller

import (
	"Moreover/internal/pkg/captcha"
	"Moreover/internal/pkg/redis"
	"Moreover/internal/pkg/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
)

var uid uuid.UUID

func GenerateCaptcha(c *gin.Context) {
	id, base64 := captcha.GenerateCaptcha()
	if id == "" {
		response.Response(c, response.ERROR, nil)
		return
	}
	uid = uuid.New()
	redis.DB.Set("captcha:"+uid.String(), id, time.Minute*5)
	response.Response(c, response.SUCCESS, gin.H{
		"id":     id,
		"base64": base64,
	})
}

func ParseCaptcha(c *gin.Context) {
	requestId := c.PostForm("captcha")
	id, err := redis.DB.Get("captcha:" + uid.String()).Result()
	if err != nil {
		fmt.Printf("parse captch from redis fail, err: %v\n", err)
		response.Response(c, response.ERROR, nil)
		return
	}
	if captcha.ParseCaptcha(id, requestId) {
		response.Response(c, response.SUCCESS, nil)
	}
}
