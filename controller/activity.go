package controller

import (
	"Moreover/internal/pkg/activity"
	"Moreover/model"
	"Moreover/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strconv"
)

func GetActivityByPage(c *gin.Context) {
	current, _ := strconv.Atoi(c.Param("current"))
	pageSize, _ := strconv.Atoi(c.Param("pageSize"))
	activities := activity.GetActivityByPage(current, pageSize)
	if activities == nil {
		response.Response(c, response.FAIL, nil)
		c.Abort()
		return
	}
	response.Response(c, response.SUCCESS, gin.H{"content": activities})
}

func PublishActivity(c *gin.Context) {
	publisher, ok := c.Get("publisher")
	if !ok {
		response.Response(c, response.AuthError, nil)
		return
	}
	var tmpActivity model.Activity
	if err := c.BindJSON(&tmpActivity); err != nil {
		response.Response(c, response.ParamError, nil)
		return
	}
	activityId := uuid.New().String()
	tmpActivity.ActivityId = activityId
	tmpActivity.Publisher = publisher.(string)
	if !activity.PublishActivity(tmpActivity) {
		response.Response(c, response.ParamError, nil)
		return
	}
	response.Response(c, response.SUCCESS, nil)
}
