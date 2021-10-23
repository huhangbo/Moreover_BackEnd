package router

import (
	"Moreover/controller"
	"Moreover/middleware/auth"
)

func activityRouter() {
	r := Router.Group("/activity")
	r.Use(auth.Auth())
	{
		r.GET("/detail/:activityId", controller.GetActivityById)

		r.POST("", controller.PublishActivity)

		r.PUT("/:activityId", controller.UpdateActivity)

		r.DELETE("/:activityId", controller.DeleteActivity)

		r.GET("/:type/:current/:pageSize", controller.GetActivityByPage)
	}
}
