package router

import (
	"Moreover/controller"
	"Moreover/middleware/auth"
)

func followRouter() {
	r := Router.Group("/follow")
	{
		r.POST("/:parentId", auth.Auth(), controller.Follow)

		r.DELETE("/:parentId", auth.Auth(), controller.UnFollow)

		r.GET("/:followType/:id/:current/:pageSize", controller.GetFollowByPage)

	}
}
