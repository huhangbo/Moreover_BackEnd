package likeController

import (
	"Moreover/internal/pkg/liked"
	"Moreover/model"
	"Moreover/pkg/response"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func PublishLike(c *gin.Context) {
	parentId := c.Param("parentId")
	likeUser := c.Param("likeUser")
	stuId, ok := c.Get("stuId")
	if !ok {
		response.Response(c, response.AuthError, nil)
		return
	}
	now := time.Now().Format("2006/01/02 15:04:05")
	tmpLike := model.Like{
		CreateTime:    now,
		UpdateTime:    now,
		LikeId:        stuId.(string) + parentId,
		ParentId:      parentId,
		LikeUser:      likeUser,
		LikePublisher: stuId.(string),
		Deleted:       0,
	}
	code := liked.PublishLike(tmpLike)
	response.Response(c, code, nil)
}

func GetLikesByPage(c *gin.Context) {
	parentId := c.Param("parentId")
	current, _ := strconv.Atoi(c.Param("current"))
	pageSize, _ := strconv.Atoi(c.Param("pageSize"))
	code, likes, page := liked.GetLikesByPage(current, pageSize, parentId)
	if code != response.SUCCESS {
		response.Response(c, code, nil)
		return
	}
	response.Response(c, code, gin.H{
		"likes": likes,
		"page":  page,
	})
}

func DeleteLike(c *gin.Context) {
	likeId := c.Param("likeId")
	stuId, ok := c.Get("stuId")
	userId := likeId[:8]
	if !ok || userId != stuId.(string) {
		response.Response(c, response.AuthError, nil)
		return
	}
	code := liked.DeleteLikeById(likeId, stuId.(string))
	response.Response(c, code, nil)
}
