package controller

import (
	"Moreover/internal/pkg/activity"
	"Moreover/model"
	"Moreover/pkg/response"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strconv"
	"time"
)

func GetActivityById(c *gin.Context) {
	activityId := c.Param("activityId")
	var activityStruct model.Activity
	code := activity.GetActivityById(activityId, &activityStruct)
	tmp, _ := json.Marshal(&activityStruct)
	activityMap := make(map[string]interface{})
	if err := json.Unmarshal(tmp, &activityMap); err != nil {
		response.Response(c, response.ERROR, nil)
	}
	response.Response(c, code, activityMap)
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
	activityId := uuid.New().String()
	tmpActivity.ActivityId = activityId
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
	activities := activity.GetActivityByPage(current, pageSize)
	if activities == nil {
		response.Response(c, response.FAIL, nil)
		c.Abort()
		return
	}
	response.Response(c, response.SUCCESS, gin.H{"content": activities})
}

func UpdateActivity(c *gin.Context) {
	stuId, ok := c.Get("stuId")
	if !ok {
		response.Response(c, response.AuthError, nil)
		return
	}
	var tmpActivity model.Activity
	if err := c.BindJSON(&tmpActivity); err != nil {
		response.Response(c, response.ParamError, nil)
		return
	}
	tmpActivity.ActivityId = c.Param("activityId")
	tmpActivity.UpdateTime = time.Now().Format("2006/01/02 15:04:05")
	code := activity.UpdateActivityById(stuId.(string), tmpActivity)
	response.Response(c, code, nil)
}

func DeleteActivity(c *gin.Context) {
	stuId, ok := c.Get("stuId")
	if !ok {
		response.Response(c, response.AuthError, nil)
		return
	}
	activityId := c.Param("activityId")
	code := activity.DeleteActivityById(stuId.(string), activityId)
	response.Response(c, code, nil)
}
