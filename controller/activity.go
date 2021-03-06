package controller

import (
	"Moreover/dao"
	"Moreover/pkg/response"
	"Moreover/service"
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
	code := service.PublishActivity(tmpActivity)
	response.Response(c, code, nil)
}

func GetActivityById(c *gin.Context) {
	stuId, _ := c.Get("stuId")
	activityId := c.Param("activityId")
	activityDetail := dao.ActivityDetailFollow{
		Activity: dao.Activity{
			ActivityId: activityId,
		},
	}
	code := service.GetActivityDetailFollow(&activityDetail, stuId.(string))
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
	code := service.UpdateActivity(tmpActivity)
	response.Response(c, code, nil)
}

func DeleteActivity(c *gin.Context) {
	stuId, _ := c.Get("stuId")
	tmpActivity := dao.Activity{
		ActivityId: c.Param("activityId"),
		Publisher:  stuId.(string),
	}
	code := service.DeleteActivity(tmpActivity)
	response.Response(c, code, nil)
}

func GetActivityByPage(c *gin.Context) {
	stuId, _ := c.Get("stuId")
	current, _ := strconv.Atoi(c.Param("current"))
	pageSize, _ := strconv.Atoi(c.Param("pageSize"))
	userId := c.Query("userId")
	switch c.Param("type") {
	case "page":
		category := c.Query("category")
		code, activities, isEnd := service.GetActivitiesByCategory(current, pageSize, stuId.(string), category)
		if code != response.SUCCESS {
			response.Response(c, code, nil)
			return
		}
		response.Response(c, code, gin.H{
			"activities": activities,
			"isEnd":      isEnd,
		})
	case "publisher":
		if userId == "" {
			userId = stuId.(string)
		}
		code, activities, isEnd := service.GetActivitiesByPublisher(current, pageSize, stuId.(string), userId)
		if code != response.SUCCESS {
			response.Response(c, code, nil)
			return
		}
		response.Response(c, code, gin.H{
			"activities": activities,
			"isEnd":      isEnd,
		})
		return
	default:
		response.Response(c, response.ParamError, nil)
	}
}
