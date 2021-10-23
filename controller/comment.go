package controller

import (
	"Moreover/dao"
	"Moreover/pkg/response"
	"Moreover/service/activity"
	"Moreover/service/comment"
	"Moreover/service/post"
	"Moreover/service/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strconv"
)

func PublishComment(c *gin.Context) {
	stuId, _ := c.Get("stuId")
	ParentId := c.Param("parentId")
	tmpComment := dao.Comment{
		CommentId: uuid.New().String(),
		ParentId:  ParentId,
		Publisher: stuId.(string),
	}
	var replier string
	if err := c.BindJSON(&tmpComment); err != nil {
		response.Response(c, response.ParamError, nil)
		return
	}
	switch c.Param("kind") {
	case "activity":
		tmpKind := dao.Activity{ActivityId: ParentId}
		code := activity.GetActivityById(&tmpKind)
		if code != response.SUCCESS {
			response.Response(c, response.ParamError, nil)
			return
		}
		replier = tmpKind.Publisher
	case "post":
		tmpKind := dao.PostDetail{Post: dao.Post{PostId: ParentId}}
		code := post.GetPostDetail(&tmpKind, stuId.(string))
		util.TopPost(tmpKind)
		if code != response.SUCCESS {
			response.Response(c, response.ParamError, nil)
			return
		}
	case "child":
		tmpKind := dao.Comment{CommentId: ParentId}
		code := comment.GetCommentById(&tmpKind)
		if code != response.SUCCESS {
			response.Response(c, response.ParamError, nil)
			return
		}
		replier = tmpKind.Publisher
	}
	tmpComment.Replier = replier
	code := comment.PublishComment(tmpComment)
	response.Response(c, code, nil)
}

func DeleteComment(c *gin.Context) {
	stuId, ok := c.Get("stuId")
	if !ok {
		response.Response(c, response.AuthError, nil)
		return
	}
	tmpComment := dao.Comment{
		CommentId: c.Param("commentId"),
	}
	code := comment.DeleteComment(tmpComment, stuId.(string))
	response.Response(c, code, nil)
}

func GetCommentsByPage(c *gin.Context) {
	current, _ := strconv.Atoi(c.Param("current"))
	pageSize, _ := strconv.Atoi(c.Param("pageSize"))
	kind := c.Param("kind")
	parentId := c.Param("parentId")
	stuId, _ := c.Get("stuId")
	code, comments, tmpPage := comment.GetCommentByIdPage(current, pageSize, parentId, stuId.(string))
	if kind == "child" {
		code, comments, tmpPage := comment.GetCommentChildrenByPage(current, pageSize, parentId, stuId.(string))
		if code != response.SUCCESS {
			response.Response(c, code, nil)
			return
		}
		response.Response(c, code, gin.H{
			"comments": comments,
			"page":     tmpPage,
		})
		return
	}
	if code != response.SUCCESS {
		response.Response(c, code, nil)
		return
	}
	response.Response(c, code, gin.H{
		"comments": comments,
		"page":     tmpPage,
	})
}
