package router

import (
	"Moreover/controller/followController"
	"Moreover/middleware/auth"
)

func followRouter() {
	r := Router.Group("/follow")
	{
		r.POST("/:follower", auth.Auth(), followController.Follow)

		r.DELETE("/:follower", auth.Auth(), followController.UnFollow)

		r.GET("/:followType/:id/:current/:pageSize", followController.GetFollowByPage)

	}
}
