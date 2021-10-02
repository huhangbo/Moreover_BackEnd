package commentController

import (
	"Moreover/internal/pkg/comment"
	"Moreover/model"
	"Moreover/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
)

func PublishComment(c *gin.Context) {
	stuId, ok := c.Get("stuId")
	id := c.Param("id")
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
		CommentId:   uuid.New().String(),
		Message:     c.PostForm("message"),
	}
	code := comment.PublishComment(tmpComment)
	response.Response(c, code, nil)
}

func GetCommentById(c *gin.Context) {
	commentId := c.Param("commentId")
	code, tmp := comment.GetCommentById(commentId)
	if code != response.SUCCESS {
		response.Response(c, code, nil)
		return
	}
	response.Response(c, code, gin.H{"content": tmp})
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
