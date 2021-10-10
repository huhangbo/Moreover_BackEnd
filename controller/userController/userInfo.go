package userController

import (
	"Moreover/internal/pkg/user"
	"Moreover/model"
	"Moreover/pkg/response"
	"github.com/gin-gonic/gin"
)

func GetUserInfoById(c *gin.Context) {
	stuId, ok := c.Get("stuId")
	if !ok {
		response.Response(c, response.AuthError, nil)
		return
	}
	userId := c.Param("userId")
	code, userInfoDetail := user.GetUserInfoDetail(userId, stuId.(string))
	if code != response.SUCCESS {
		response.Response(c, code, nil)
		return
	}
	response.Response(c, code, gin.H{"userInfo": userInfoDetail})
}

func UpdateAvatar(c *gin.Context) {
	userId, ok := c.Get("stuId")
	if !ok {
		response.Response(c, response.AuthError, nil)
		return
	}
	var tmpUserInfo model.UserInfo
	if err := c.BindJSON(&tmpUserInfo); err != nil {
		response.Response(c, response.ERROR, nil)
	}
	if tmpUserInfo.Avatar == "" {
		response.Response(c, response.ParamError, nil)
		return
	}
	code := user.UpdateUserAvatar(tmpUserInfo.Avatar, userId.(string))
	response.Response(c, code, nil)
}

func UpdateSex(c *gin.Context) {
	userId, ok := c.Get("stuId")
	if !ok {
		response.Response(c, response.AuthError, nil)
		return
	}
	var tmpUserInfo model.UserInfo
	if err := c.BindJSON(&tmpUserInfo); err != nil {
		response.Response(c, response.ERROR, nil)
	}
	if tmpUserInfo.Sex == "" {
		response.Response(c, response.ParamError, nil)
		return
	}
	code := user.UpdateUserSex(tmpUserInfo.Sex, userId.(string))
	response.Response(c, code, nil)
}

func UpdateNickname(c *gin.Context) {
	userId, ok := c.Get("stuId")
	if !ok {
		response.Response(c, response.AuthError, nil)
		return
	}
	var tmpUserInfo model.UserInfo
	if err := c.BindJSON(&tmpUserInfo); err != nil {
		response.Response(c, response.ERROR, nil)
	}
	if tmpUserInfo.Nickname == "" {
		response.Response(c, response.ParamError, nil)
		return
	}
	code := user.UpdateUserNickname(tmpUserInfo.Nickname, userId.(string))
	response.Response(c, code, nil)
}

func UpdateDescription(c *gin.Context) {
	userId, ok := c.Get("stuId")
	if !ok {
		response.Response(c, response.AuthError, nil)
		return
	}
	var tmpUserInfo model.UserInfo
	if err := c.BindJSON(&tmpUserInfo); err != nil {
		response.Response(c, response.ERROR, nil)
	}
	if tmpUserInfo.Description == "" {
		response.Response(c, response.ParamError, nil)
		return
	}
	code := user.UpdateUserDescription(tmpUserInfo.Description, userId.(string))
	response.Response(c, code, nil)
}
