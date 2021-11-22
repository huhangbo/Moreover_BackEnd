package router

import (
	"Moreover/controller"
	"Moreover/middleware/auth"
)

func commentRouter() {
	r := Router.Group("/comment")
	r.Use(auth.Auth())
	{
		r.POST("/:kind/:parentId", controller.PublishComment)

		r.GET("/:kind/:parentId/:current/:pageSize", controller.GetCommentsByPage)

		r.DELETE("/:commentId", controller.DeleteComment)
	}
}
