package controller

import (
	"Moreover/dao"
	"Moreover/pkg/response"
	"Moreover/service/follow"
	"github.com/gin-gonic/gin"
	"strconv"
)

func Follow(c *gin.Context) {
	stuId, _ := c.Get("stuId")
	parentId := c.Param("parentId")
	tmpFollow := dao.Follow{
		Parent:    parentId,
		Publisher: stuId.(string),
	}
	code := follow.PublishFollow(tmpFollow)
	response.Response(c, code, nil)
}

func UnFollow(c *gin.Context) {
	stuId, _ := c.Get("stuId")
	parentId := c.Param("parentId")
	tmpFollow := dao.Follow{
		Parent:    parentId,
		Publisher: stuId.(string),
	}
	code := follow.Unfollow(tmpFollow)
	response.Response(c, code, nil)
}

func GetFollowByPage(c *gin.Context) {
	current, _ := strconv.Atoi(c.Param("current"))
	pageSize, _ := strconv.Atoi(c.Param("pageSize"))
	followType := c.Param("followType")
	id := c.Param("id")
	var tmp string
	switch followType {
	case "publisher":
		tmp = "parent"
	case "parent":
		tmp = "publisher"
	default:
		response.Response(c, response.ParamError, nil)
		return
	}
	code, follows, tmpPage := follow.GetFollowById(current, pageSize, id, followType, tmp)
	if code != response.SUCCESS {
		response.Response(c, code, nil)
		return
	}
	response.Response(c, code, gin.H{"followers": follows, "page": tmpPage})
}