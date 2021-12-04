package router

import (
	"Moreover/controller"
	"Moreover/middleware"
)

func likeRouter() {
	r := Router.Group("/liked")
	r.Use(middleware.Auth())
	{
		r.POST("/:kind/:parentId", controller.PublishLike)

		r.GET("/:parentId/:current/:pageSize", controller.GetLikesByPage)

		r.DELETE("/:kind/:parentId", controller.DeleteLike)
	}
}
