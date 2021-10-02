package router

import (
	"Moreover/controller/commentController"
	"Moreover/middleware/auth"
)

func CommentRouter() {
	r := Router.Group("/comment")
	{
		r.POST("/:parentId", auth.Auth(), commentController.PublishComment)

		r.GET("/:parentId/:current/:pageSize", commentController.GetCommentsByPage)

		r.DELETE("/:commentId", auth.Auth(), commentController.DeleteComment)

	}
}
