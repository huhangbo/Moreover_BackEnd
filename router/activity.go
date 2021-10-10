package router

import (
	"Moreover/controller/activityController"
	"Moreover/middleware/auth"
)

func activityRouter() {
	r := Router.Group("/activity")
	r.Use(auth.Auth())
	{
		r.GET("/:activityId", activityController.GetActivityById)

		r.POST("", activityController.PublishActivity)

		r.PUT("/:activityId", activityController.UpdateActivity)

		r.DELETE("/:activityId", activityController.DeleteActivity)

		r.GET("/page/:current/:pageSize", activityController.GetActivityByPage)

		r.GET("/publish/:current/:pageSize", activityController.GetActivitiesByPublisher)

	}
}
