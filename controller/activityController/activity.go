package activityController

import (
	"Moreover/internal/pkg/activity"
	"Moreover/model"
	"Moreover/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strconv"
	"time"
)

func GetActivityById(c *gin.Context) {
	activityId := c.Param("activityId")
	code, tmpActivity := activity.GetActivityById(activityId)
	if code != response.SUCCESS {
		response.Response(c, code, nil)
		return
	}
	response.Response(c, code, gin.H{"content": tmpActivity})
}

func PublishActivity(c *gin.Context) {
	publisher, ok := c.Get("stuId")
	if !ok {
		response.Response(c, response.AuthError, nil)
		return
	}
	var tmpActivity model.Activity
	if err := c.BindJSON(&tmpActivity); err != nil {
		response.Response(c, response.ParamError, nil)
		return
	}
	tmpActivity.ActivityId = uuid.New().String()
	tmpActivity.Publisher = publisher.(string)
	now := time.Now().Format("2006/01/02 15:04:05")
	tmpActivity.PublishTime = now
	tmpActivity.UpdateTime = now
	code := activity.PublishActivity(tmpActivity)
	response.Response(c, code, nil)
}

func GetActivityByPage(c *gin.Context) {
	current, _ := strconv.Atoi(c.Param("current"))
	pageSize, _ := strconv.Atoi(c.Param("pageSize"))
	category := c.Query("category")
	code, activities, page := activity.GetActivitiesByPade(current, pageSize, category)
	if code != response.SUCCESS {
		response.Response(c, code, nil)
		return
	}
	response.Response(c, code, gin.H{
		"activities": activities,
		"page":       page,
	})
}

func UpdateActivity(c *gin.Context) {
	stuId, ok := c.Get("stuId")
	if !ok {
		response.Response(c, response.AuthError, nil)
		return
	}
	activityId := c.Param("activityId")
	var tmpActivity model.Activity
	if err := c.BindJSON(&tmpActivity); err != nil {
		response.Response(c, response.ParamError, nil)
		return
	}
	codeOld, oldActivity := activity.GetActivityById(activityId)
	if codeOld != response.SUCCESS {
		response.Response(c, codeOld, nil)
		return
	}
	tmpActivity.ActivityId = activityId
	tmpActivity.UpdateTime = time.Now().Format("2006/01/02 15:04:05")
	tmpActivity.PublishTime = time.Now().Format("2006/01/02 15:04:05")
	tmpActivity.Publisher = stuId.(string)
	codeUpdate := activity.UpdateActivityById(tmpActivity, oldActivity)
	response.Response(c, codeUpdate, nil)
}

func DeleteActivity(c *gin.Context) {
	stuId, ok := c.Get("stuId")
	if !ok {
		response.Response(c, response.AuthError, nil)
		return
	}
	activityId := c.Param("activityId")
	code := activity.DeleteActivityById(activityId, stuId.(string))
	response.Response(c, code, nil)
}
