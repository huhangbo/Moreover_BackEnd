package controller

import (
	"Moreover/dao"
	"Moreover/pkg/response"
	"Moreover/service"
	"Moreover/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strconv"
	"time"
)

func PublishComment(c *gin.Context) {
	stuId, _ := c.Get("stuId")
	ParentId := c.Param("parentId")
	kind := c.Param("kind")
	tmpComment := dao.Comment{
		CommentId: uuid.New().String(),
		ParentId:  ParentId,
		Publisher: stuId.(string),
	}
	tmpMessage := dao.Message{
		CreatedAt: time.Now(),
		Parent:    tmpComment.CommentId,
		Action:    "comment",
		Kind:      kind,
		Publisher: stuId.(string),
	}
	var replier string
	if err := c.BindJSON(&tmpComment); err != nil {
		response.Response(c, response.ParamError, nil)
		return
	}
	switch kind {
	case "activity":
		tmpKind := dao.Activity{ActivityId: ParentId}
		code := service.GetActivityById(&tmpKind)
		if code != response.SUCCESS {
			response.Response(c, response.ParamError, nil)
			return
		}
		replier = tmpKind.Publisher
	case "post":
		tmpKind := dao.PostDetail{Post: dao.Post{PostId: ParentId}}
		code := service.GetPostDetail(&tmpKind, stuId.(string))
		if code != response.SUCCESS {
			response.Response(c, response.ParamError, nil)
			return
		}
		if err := util.TopPost(ParentId, "comment"); err != nil {
			return
		}
	case "child":
		tmpKind := dao.Comment{CommentId: ParentId}
		code := service.GetCommentById(&tmpKind)
		if code != response.SUCCESS {
			response.Response(c, response.ParamError, nil)
			return
		}
		replier = tmpKind.Publisher
	default:
		response.Response(c, response.ParamError, nil)
		return
	}
	tmpComment.Replier = replier
	tmpMessage.Receiver = replier
	tmpMessage.Detail = tmpComment.Message
	if err := service.PublishMessage(tmpMessage); err == nil {
		service.UserMap.PostMessage(&tmpMessage)
	}
	code := service.PublishComment(tmpComment)
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
	code := service.DeleteComment(tmpComment, stuId.(string))
	response.Response(c, code, nil)
}

func GetCommentsByPage(c *gin.Context) {
	current, _ := strconv.Atoi(c.Param("current"))
	pageSize, _ := strconv.Atoi(c.Param("pageSize"))
	kind := c.Param("kind")
	parentId := c.Param("parentId")
	stuId, _ := c.Get("stuId")
	switch kind {
	case "parent":
		{
			code, comments, tmpPage := service.GetCommentByIdPage(current, pageSize, parentId, stuId.(string))
			if code != response.SUCCESS {
				response.Response(c, code, nil)
				return
			}
			response.Response(c, code, gin.H{
				"comments": comments,
				"page":     tmpPage,
			})
		}
	case "child":
		{
			code, comments, tmpPage := service.GetCommentChildrenByPage(current, pageSize, parentId, stuId.(string))
			if code != response.SUCCESS {
				response.Response(c, code, nil)
				return
			}
			response.Response(c, code, gin.H{
				"comments": comments,
				"page":     tmpPage,
			})
		}
	default:
		response.Response(c, response.ParamError, nil)
	}
}
