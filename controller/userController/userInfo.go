package userController

import (
	"Moreover/internal/pkg/user"
	"Moreover/model"
	"Moreover/pkg/response"
	"github.com/gin-gonic/gin"
)

func GetUserInfoById(c *gin.Context) {
	userId := c.Param("userId")
	code, userInfo := user.GetUserInfo(userId)
	if code != response.SUCCESS {
		response.Response(c, code, nil)
		return
	}
	response.Response(c, code, gin.H{"userInfo": userInfo})
}

func UpdateAvatar(c *gin.Context) {
	userId := c.Param("userId")
	var tmpUserInfo model.UserInfo
	if err := c.BindJSON(&tmpUserInfo); err != nil {
		response.Response(c, response.ERROR, nil)
	}
	if tmpUserInfo.Avatar == "" {
		response.Response(c, response.ParamError, nil)
		return
	}
	code := user.UpdateUserAvatar(tmpUserInfo.Avatar, userId)
	response.Response(c, code, nil)
}

func UpdateSex(c *gin.Context) {
	userId := c.Param("userId")
	var tmpUserInfo model.UserInfo
	if err := c.BindJSON(&tmpUserInfo); err != nil {
		response.Response(c, response.ERROR, nil)
	}
	if tmpUserInfo.Sex == "" {
		response.Response(c, response.ParamError, nil)
		return
	}
	code := user.UpdateUserSex(tmpUserInfo.Sex, userId)
	response.Response(c, code, nil)
}

func UpdateNickname(c *gin.Context) {
	userId := c.Param("userId")
	var tmpUserInfo model.UserInfo
	if err := c.BindJSON(&tmpUserInfo); err != nil {
		response.Response(c, response.ERROR, nil)
	}
	if tmpUserInfo.Nickname == "" {
		response.Response(c, response.ParamError, nil)
		return
	}
	code := user.UpdateUserNickname(tmpUserInfo.Nickname, userId)
	response.Response(c, code, nil)
}

func UpdateDescription(c *gin.Context) {
	userId := c.Param("userId")
	var tmpUserInfo model.UserInfo
	if err := c.BindJSON(&tmpUserInfo); err != nil {
		response.Response(c, response.ERROR, nil)
	}
	if tmpUserInfo.Description == "" {
		response.Response(c, response.ParamError, nil)
		return
	}
	code := user.UpdateUserDescription(tmpUserInfo.Description, userId)
	response.Response(c, code, nil)
}
