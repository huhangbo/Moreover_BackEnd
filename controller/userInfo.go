package controller

import (
	"Moreover/dao"
	"Moreover/pkg/response"
	"Moreover/service"
	"github.com/gin-gonic/gin"
)

func GetUserInfoById(c *gin.Context) {
	stuId, _ := c.Get("stuId")
	userId := c.Param("userId")
	tmpUserDetail := dao.UserInfoDetail{
		UserInfo: dao.UserInfo{
			StudentId: userId,
		},
	}
	code := service.GetUserInfoDetail(&tmpUserDetail, stuId.(string))
	if code != response.SUCCESS {
		response.Response(c, code, nil)
		return
	}
	response.Response(c, code, gin.H{"userInfo": tmpUserDetail})
}

func UpdateUserInfo(c *gin.Context) {
	stuId, _ := c.Get("stuId")
	tmpUserInfo := dao.UserInfo{
		StudentId: stuId.(string),
	}
	if err := c.BindJSON(&tmpUserInfo); err != nil {
		response.Response(c, response.ParamError, nil)
	}
	code := service.UpdateUserInfo(tmpUserInfo)
	response.Response(c, code, nil)
}
