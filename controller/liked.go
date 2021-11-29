package controller

import (
	"Moreover/dao"
	"Moreover/pkg/response"
	"Moreover/service"
	"Moreover/util"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func PublishLike(c *gin.Context) {
	parentId := c.Param("parentId")
	stuId, _ := c.Get("stuId")
	kind := c.Param("kind")
	tmpLike := dao.Liked{
		Parent:    parentId,
		Publisher: stuId.(string),
	}
	tmpMessage := dao.Message{
		CreatedAt: time.Now(),
		Parent:    parentId,
		Action:    "liked",
		Kind:      kind,
		Publisher: stuId.(string),
	}
	switch kind {
	case "activity":
		tmp := dao.Activity{ActivityId: parentId}
		if code := service.GetActivityById(&tmp); code != response.SUCCESS {
			response.Response(c, response.ParamError, nil)
			return
		}
		tmpLike.Liker = tmp.Publisher
		tmpMessage.Receiver = tmp.Publisher
	case "comment":
		tmp := dao.Comment{CommentId: parentId}
		if code := service.GetCommentById(&tmp); code != response.SUCCESS {
			response.Response(c, response.ParamError, nil)
			return
		}
		tmpLike.Liker = tmp.Publisher
		tmpMessage.Receiver = tmp.Publisher
	case "post":
		tmp := dao.PostDetail{Post: dao.Post{PostId: parentId}}
		if code := service.GetPostDetail(&tmp, stuId.(string)); code != response.SUCCESS {
			response.Response(c, response.ParamError, nil)
			return
		}
		tmpLike.Liker = tmp.Publisher
		tmpMessage.Receiver = tmp.Publisher
	default:
		response.Response(c, response.ParamError, nil)
		return
	}
	code := service.PublishLike(tmpLike)
	if code == response.SUCCESS {
		if err := service.PublishMessage(tmpMessage); err == nil {
			service.UserMap.PostMessage(&tmpMessage)
		}
		_ = util.TopPost(parentId, "liked")
	}
	response.Response(c, code, nil)
}

func GetLikesByPage(c *gin.Context) {
	parentId := c.Param("parentId")
	current, _ := strconv.Atoi(c.Param("current"))
	pageSize, _ := strconv.Atoi(c.Param("pageSize"))
	code, likes, isEnd := service.GetLikeByPage(current, pageSize, parentId)
	if code != response.SUCCESS {
		response.Response(c, code, nil)
		return
	}
	response.Response(c, code, gin.H{
		"likes": likes,
		"isEnd": isEnd,
	})
}

func DeleteLike(c *gin.Context) {
	parentId := c.Param("parentId")
	stuId, _ := c.Get("stuId")
	tmpLiked := dao.Liked{
		Parent:    parentId,
		Publisher: stuId.(string),
	}
	code := service.UnLike(tmpLiked)
	if code == response.SUCCESS && c.Param("kind") == "post" {
		_ = util.TopPost(parentId, "dislike")
	}
	response.Response(c, code, nil)
}
