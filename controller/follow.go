package controller

import (
	"Moreover/model"
	"Moreover/pkg/response"
	follow2 "Moreover/service/follow"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func Follow(c *gin.Context) {
	stuId, ok := c.Get("stuId")
	if !ok {
		response.Response(c, response.AuthError, nil)
		return
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	follower := c.Param("follower")
	tmpFollow := model.Follow{
		CreateTime: now,
		UpdateTime: now,
		Fan:        stuId.(string),
		Follower:   follower,
	}
	code := follow2.PublishFollow(tmpFollow)
	response.Response(c, code, nil)
}

func UnFollow(c *gin.Context) {
	stuId, ok := c.Get("stuId")
	if !ok {
		response.Response(c, response.AuthError, nil)
		return
	}
	follower := c.Param("follower")
	code := follow2.Unfollow(follower, stuId.(string))
	response.Response(c, code, nil)
}

func GetFollowByPage(c *gin.Context) {
	current, _ := strconv.Atoi(c.Param("current"))
	pageSize, _ := strconv.Atoi(c.Param("pageSize"))
	followType := c.Param("followType")
	id := c.Param("id")
	code, follows, tmpPage := follow2.GetFollowById(current, pageSize, id, followType)
	if code != response.SUCCESS {
		response.Response(c, code, nil)
		return
	}
	response.Response(c, code, gin.H{"followers": follows, "page": tmpPage})
}
