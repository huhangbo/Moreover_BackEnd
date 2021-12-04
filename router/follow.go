package router

import (
	"Moreover/controller"
	"Moreover/middleware"
)

func followRouter() {
	r := Router.Group("/follow")
	{
		r.POST("/:parentId", middleware.Auth(), controller.Follow)

		r.DELETE("/:parentId", middleware.Auth(), controller.UnFollow)

		r.GET("/:followType/:id/:current/:pageSize", controller.GetFollowByPage)
	}
}
