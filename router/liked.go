package router

import (
	"Moreover/controller/likeController"
	"Moreover/middleware/auth"
)

func likeRouter() {
	r := Router.Group("/liked")
	{
		r.POST("/:parentId", auth.Auth(), likeController.PublishLike)

		r.GET("/:parentId/:current/:pageSize", likeController.GetLikesByPage)

		r.DELETE("/:parentId", auth.Auth(), likeController.DeleteLike)
	}
}
