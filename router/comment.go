package router

import (
	"Moreover/controller/commentController"
	"Moreover/middleware/auth"
)

func CommentRouter() {
	r := Router.Group("/comment")
	{
		r.POST("/:id", auth.Auth(), commentController.PublishComment)

		r.GET("/:commentId", auth.Auth(), commentController.GetCommentById)

		r.DELETE("/:commentId", auth.Auth(), commentController.DeleteComment)
	}
}
