package router

import (
	"Moreover/controller/activityController"
	"Moreover/middleware/auth"
)

func activityRouter() {
	r := Router.Group("/activity")
	{
		r.GET("/:activityId", activityController.GetActivityById)

		r.POST("", auth.Auth(), activityController.PublishActivity)

		r.PUT("/:activityId", auth.Auth(), activityController.UpdateActivity)

		r.DELETE("/:activityId", auth.Auth(), activityController.DeleteActivity)

		r.GET("/page/:current/:pageSize", activityController.GetActivityByPage)

	}
}
