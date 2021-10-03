package commentController

import (
	"Moreover/internal/pkg/comment"
	"Moreover/model"
	"Moreover/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strconv"
	"time"
)

func PublishComment(c *gin.Context) {
	id := c.Param("parentId")
	stuId, ok := c.Get("stuId")
	if !ok {
		response.Response(c, response.AuthError, nil)
		return
	}
	now := time.Now().Format("2006/01/02 15:04:05")
	var tmpComment = model.Comment{
		PublishTime: now,
		UpdateTime:  now,
		ParentID:    id,
		Publisher:   stuId.(string),
		Replier:     c.PostForm("replier"),
		CommentId:   uuid.New().String(),
		Message:     c.PostForm("message"),
	}
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
	childCurrent, _ := strconv.Atoi(c.Query("childCurrent"))
	childSize, _ := strconv.Atoi(c.Query("childSize"))
	parentId := c.Param("parentId")
	code, comments, tmpPage := comment.GetCommentsByIdPage(current, pageSize, parentId)
	if code != response.SUCCESS {
		response.Response(c, code, nil)
		return
	}
	for i := 0; i < len(comments); i++ {
		code, childComment, childPage := comment.GetCommentsByIdPage(childCurrent, childSize, comments[i].CommentId)
		if code != response.SUCCESS {
			response.Response(c, code, nil)
			return
		}
		if len(childComment) == 0 {
			continue
		}
		comments[i].Children.Comments = append(comments[i].Children.Comments, childComment...)
		comments[i].Children.Page = childPage
	}
	response.Response(c, code, gin.H{
		"comments": comments,
		"page":     tmpPage,
	})
}
