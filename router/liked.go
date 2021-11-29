package router

import (
	"Moreover/controller"
	"Moreover/middleware"
)

func likeRouter() {
	r := Router.Group("/liked")
	{
		r.POST("/:kind/:parentId", middleware.Auth(), controller.PublishLike)

		r.GET("/:parentId/:current/:pageSize", controller.GetLikesByPage)

		r.DELETE("/:kind/:parentId", middleware.Auth(), controller.DeleteLike)
	}
}
