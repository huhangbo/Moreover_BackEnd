package router

import (
	"Moreover/controller"
	"Moreover/middleware"
)

func messageRouter() {
	r := Router.Group("/message").Use(middleware.Auth())
	{
		r.GET("/connect", controller.HandleSSE)

		r.GET("/:action/:current/:pageSize", controller.GetMessages)
	}
}
