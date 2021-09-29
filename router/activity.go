package router

import (
	"Moreover/controller"
	"Moreover/middleware/auth"
)

func ActivityRouter() {
	r := Router.Group("/activity")
	{
		r.GET("/:activityId", controller.GetActivityById)

		r.POST("", auth.Auth(), controller.PublishActivity)

		r.PUT("/:activityId", auth.Auth(), controller.UpdateActivity)

		r.DELETE("/:activityId", auth.Auth(), controller.DeleteActivity)

		r.GET("/page/:current/:pageSize", controller.GetActivityByPage)

	}
}
