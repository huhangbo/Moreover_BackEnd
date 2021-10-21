package router

import (
	"Moreover/controller"
	"Moreover/middleware/auth"
)

func likeRouter() {
	r := Router.Group("/liked")
	{
		r.POST("/:kind/:parentId", auth.Auth(), controller.PublishLike)

		r.GET("/:parentId/:current/:pageSize", controller.GetLikesByPage)

		r.DELETE("/:kind/:parentId", auth.Auth(), controller.DeleteLike)
	}
}
