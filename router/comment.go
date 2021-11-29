package router

import (
	"Moreover/controller"
	"Moreover/middleware"
)

func commentRouter() {
	r := Router.Group("/comment")
	r.Use(middleware.Auth())
	{
		r.POST("/:kind/:parentId", controller.PublishComment)

		r.GET("/:kind/:parentId/:current/:pageSize", controller.GetCommentsByPage)

		r.DELETE("/:commentId", controller.DeleteComment)
	}
}
