package controller

import (
	"Moreover/dao"
	"Moreover/pkg/response"
	"Moreover/service/activity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strconv"
	"time"
)

func PublishActivity(c *gin.Context) {
	stuId, _ := c.Get("stuId")
	tmpActivity := dao.Activity{
		CreatedAt:  time.Now().Local(),
		ActivityId: uuid.New().String(),
		Publisher:  stuId.(string),
	}
	if err := c.BindJSON(&tmpActivity); err != nil {
		response.Response(c, response.ParamError, nil)
		return
	}
	code := activity.PublishActivity(tmpActivity)
	response.Response(c, code, nil)
}

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

func UpdateActivity(c *gin.Context) {
	stuId, _ := c.Get("stuId")
	tmpActivity := dao.Activity{
		Publisher: stuId.(string),
	}
	if err := c.BindJSON(&tmpActivity); err != nil {
		response.Response(c, response.ParamError, nil)
		return
	}
	tmpActivity.ActivityId = c.Param("activityId")
	code := activity.UpdateActivity(tmpActivity)
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

func GetActivityByPage(c *gin.Context) {
	stuId, _ := c.Get("stuId")
	current, _ := strconv.ParseInt(c.Param("current"), 10, 64)
	pageSize, _ := strconv.ParseInt(c.Param("pageSize"), 10, 64)
	switch c.Param("type") {
	case "page":
		category := c.Query("category")
		code, activities, page := activity.GetActivitiesByPade(current, pageSize, stuId.(string), category)
		if code != response.SUCCESS {
			response.Response(c, code, nil)
			return
		}
		response.Response(c, code, gin.H{
			"activities": activities,
			"page":       page,
		})
	case "publisher":
		code, activities, page := activity.GetActivitiesByPublisher(current, pageSize, stuId.(string))
		if code != response.SUCCESS {
			response.Response(c, code, nil)
			return
		}
		response.Response(c, code, gin.H{
			"activities": activities,
			"page":       page,
		})
		return
	default:
		response.Response(c, response.ParamError, nil)
	}
}
