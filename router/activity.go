package router

import (
	"Moreover/controller"
	"Moreover/middleware/auth"
)

func ActivityRouter() {
	r := Router.Group("/activity")
	{
		r.GET("/:current/:pageSize", controller.GetActivityByPage)

		r.PUT("", auth.Auth(), controller.PublishActivity)
	}
}
