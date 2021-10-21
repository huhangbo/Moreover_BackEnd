package controller

import (
	"Moreover/dao"
	"Moreover/pkg/response"
	"Moreover/service/activity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strconv"
)

func GetActivityById(c *gin.Context) {
	stuId, _ := c.Get("stuId")
	activityId := c.Param("activityId")
	activityDetail := dao.ActivityDetail{
		Activity: dao.Activity{
			ActivityId: activityId,
		},
	}
	code := activity.GetActivityDetailById(&activityDetail, stuId.(string))
	if code != response.SUCCESS {
		response.Response(c, code, nil)
		return
	}
	response.Response(c, code, gin.H{"content": activityDetail})
}

func PublishActivity(c *gin.Context) {
	publisher, ok := c.Get("stuId")
	if !ok {
		response.Response(c, response.AuthError, nil)
		return
	}
	var tmpActivity dao.Activity
	if err := c.BindJSON(&tmpActivity); err != nil {
		response.Response(c, response.ParamError, nil)
		return
	}
	tmpActivity.ActivityId = uuid.New().String()
	tmpActivity.Publisher = publisher.(string)
	code := activity.PublishActivity(tmpActivity)
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
	code, activities, page := activity.GetActivitiesByPade(current, pageSize, category, stuId.(string))
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
	code, tmpActivities, tmpPage := activity.GetActivitiesByPublisher(current, pageSize, stuId.(string))
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
	var tmpActivity dao.Activity
	if err := c.BindJSON(&tmpActivity); err != nil {
		response.Response(c, response.ParamError, nil)
		return
	}
	tmpActivity.ActivityId = c.Param("activityId")
	code := activity.UpdateActivity(tmpActivity, stuId.(string))
	response.Response(c, code, nil)
}

func DeleteActivity(c *gin.Context) {
	stuId, _ := c.Get("stuId")
	tmpActivity := dao.Activity{
		ActivityId: c.Param("activityId"),
	}
	code := activity.DeleteActivity(tmpActivity, stuId.(string))
	response.Response(c, code, nil)
}
