package commentController

import (
	"Moreover/internal/pkg/activity"
	"Moreover/internal/pkg/comment"
	"Moreover/model"
	"Moreover/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strconv"
	"time"
)

func PublishComment(c *gin.Context) {
	stuId, ok := c.Get("stuId")
	if !ok {
		response.Response(c, response.AuthError, nil)
		return
	}
	ParentId := c.Param("parentId")
	var tmpComment model.Comment
	var replier string
	if err := c.BindJSON(&tmpComment); err != nil {
		response.Response(c, response.ParamError, nil)
		return
	}
	switch c.Param("kind") {
	case "activity":
		code, tmpKind := activity.GetActivityById(ParentId)
		if code != response.SUCCESS {
			response.Response(c, response.ParamError, nil)
			return
		}
		replier = tmpKind.Publisher
	case "comment":
		code, tmpKind := comment.GetCommentById(ParentId)
		if code != response.SUCCESS {
			response.Response(c, response.ParamError, nil)
			return
		}
		replier = tmpKind.Publisher
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	tmpComment.CreateTime = now
	tmpComment.UpdateTime = now
	tmpComment.ParentID = ParentId
	tmpComment.Publisher = stuId.(string)
	tmpComment.Replier = replier
	tmpComment.CommentId = uuid.New().String()
	code := comment.PublishComment(tmpComment)
	response.Response(c, code, nil)
}

func DeleteComment(c *gin.Context) {
	stuId, ok := c.Get("stuId")
	if !ok {
		response.Response(c, response.AuthError, nil)
		return
	}
	commentId := c.Param("commentId")
	code := comment.DeleteCommentById(commentId, stuId.(string))
	response.Response(c, code, nil)
}

func GetCommentsByPage(c *gin.Context) {
	current, _ := strconv.Atoi(c.Param("current"))
	pageSize, _ := strconv.Atoi(c.Param("pageSize"))
	childSize, _ := strconv.Atoi(c.Query("childSize"))
	parentId := c.Param("parentId")
	stuId, ok := c.Get("stuId")
	if !ok {
		response.Response(c, response.AuthError, nil)
		return
	}
	code, comments, tmpPage := comment.GetCommentByIdPage(current, pageSize, parentId, stuId.(string))
	if code != response.SUCCESS {
		response.Response(c, code, nil)
		return
	}
	var parentComments []model.ParentComment
	for i := 0; i < len(comments); i++ {
		code, childComment := comment.GetPreChildCById(childSize, comments[i].CommentId, stuId.(string))
		if code != response.SUCCESS {
			response.Response(c, code, nil)
			return
		}
		parentComments = append(parentComments, model.ParentComment{
			CommentDetail: comments[i],
			Children:      childComment,
		})
	}
	response.Response(c, code, gin.H{
		"comments": parentComments,
		"page":     tmpPage,
	})
}
