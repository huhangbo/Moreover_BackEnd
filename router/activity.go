package router

import (
	"Moreover/controller"
	"Moreover/middleware"
)

func activityRouter() {
	r := Router.Group("/activity")
	r.Use(middleware.Auth())
	{
		r.GET("/detail/:activityId", controller.GetActivityById)

		r.POST("", controller.PublishActivity)

		r.PUT("/:activityId", controller.UpdateActivity)

		r.DELETE("/:activityId", controller.DeleteActivity)

		r.GET("/:type/:current/:pageSize", controller.GetActivityByPage)
	}
}
