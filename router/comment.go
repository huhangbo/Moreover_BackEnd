package router

import (
	"Moreover/controller/commentController"
	"Moreover/middleware/auth"
)

func commentRouter() {
	r := Router.Group("/comment")
	r.Use(auth.Auth())
	{
		r.POST("/:kind/:parentId", commentController.PublishComment)

		r.GET("/:parentId/:current/:pageSize", commentController.GetCommentsByPage)

		r.DELETE("/:commentId", commentController.DeleteComment)

	}
}
