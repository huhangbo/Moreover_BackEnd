package controller

import (
	"Moreover/dao"
	"Moreover/pkg/response"
	"Moreover/service/activity"
	"Moreover/service/comment"
	"Moreover/service/liked"
	"github.com/gin-gonic/gin"
	"strconv"
)

func PublishLike(c *gin.Context) {
	parentId := c.Param("parentId")
	stuId, _ := c.Get("stuId")
	kind := c.Param("kind")
	var likeUser string
	switch kind {
	case "activity":
		tmp := dao.Activity{
			ActivityId: parentId,
		}
		activity.GetActivityById(&tmp)
		likeUser = tmp.Publisher
	case "comment":
		tmp := dao.Comment{
			CommentId: parentId,
		}
		comment.GetCommentById(&tmp)
		likeUser = tmp.Publisher
	default:
		response.Response(c, response.ParamError, nil)
		return
	}
	tmpLike := dao.Liked{
		ParentId:  parentId,
		LikeUser:  likeUser,
		Publisher: stuId.(string),
	}
	code := liked.PublishLike(tmpLike)
	response.Response(c, code, nil)
}

func GetLikesByPage(c *gin.Context) {
	parentId := c.Param("parentId")
	current, _ := strconv.Atoi(c.Param("current"))
	pageSize, _ := strconv.Atoi(c.Param("pageSize"))
	code, likes, page := liked.GetLikeByPage(current, pageSize, parentId)
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
	parentId := c.Param("parentId")
	stuId, _ := c.Get("stuId")
	tmpLiked := dao.Liked{
		ParentId:  parentId,
		Publisher: stuId.(string),
	}
	code := liked.UnLike(tmpLiked)
	response.Response(c, code, nil)
}