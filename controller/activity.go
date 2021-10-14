package controller

import (
	"Moreover/model"
	"Moreover/pkg/response"
	activity2 "Moreover/service/activity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strconv"
	"time"
)

func GetActivityById(c *gin.Context) {
	stuId, ok := c.Get("stuId")
	if !ok {
		response.Response(c, response.AuthError, nil)
		return
	}
	activityId := c.Param("activityId")
	code, tmpActivity := activity2.GetActivityDetailById(activityId, stuId.(string))
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
	now := time.Now().Format("2006-01-02 15:04:05")
	tmpActivity.CreateTime = now
	tmpActivity.UpdateTime = now
	code := activity2.PublishActivity(tmpActivity)
	response.Response(c, code, nil)
}

func GetActivityByPage(c *gin.Context) {
	stuId, ok := c.Get("stuId")
	if !ok {
		response.Response(c, response.AuthError, nil)
		return
	}
	current, _ := strconv.Atoi(c.Param("current"))
	pageSize, _ := strconv.Atoi(c.Param("pageSize"))
	category := c.Query("category")
	code, activities, page := activity2.GetActivitiesByPade(current, pageSize, category, stuId.(string))
	if code != response.SUCCESS {
		response.Response(c, code, nil)
		return
	}
	response.Response(c, code, gin.H{
		"activities": activities,
		"page":       page,
	})
}

func GetActivitiesByPublisher(c *gin.Context) {
	stuId, ok := c.Get("stuId")
	if !ok {
		response.Response(c, response.AuthError, nil)
		return
	}
	current, _ := strconv.Atoi(c.Param("current"))
	pageSize, _ := strconv.Atoi(c.Param("pageSize"))
	code, tmpActivities, tmpPage := activity2.GetActivityPublishedFromMysql(current, pageSize, stuId.(string))
	if code != response.SUCCESS {
		response.Response(c, code, nil)
		return
	}
	response.Response(c, code, gin.H{
		"activities": tmpActivities,
		"page":       tmpPage,
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
	code, oldActivity := activity2.GetActivityById(activityId)
	if code != response.SUCCESS {
		response.Response(c, code, nil)
		return
	}
	if oldActivity.Publisher != stuId.(string) {
		response.Response(c, response.AuthError, nil)
		return
	}
	tmpActivity.ActivityId = activityId
	tmpActivity.CreateTime = oldActivity.CreateTime
	tmpActivity.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	tmpActivity.Publisher = stuId.(string)
	code = activity2.UpdateActivityById(tmpActivity, oldActivity)
	response.Response(c, code, nil)
}

func DeleteActivity(c *gin.Context) {
	stuId, ok := c.Get("stuId")
	if !ok {
		response.Response(c, response.AuthError, nil)
		return
	}
	activityId := c.Param("activityId")
	code := activity2.DeleteActivityById(activityId, stuId.(string))
	response.Response(c, code, nil)
}
