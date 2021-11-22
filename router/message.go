package router

import (
	"Moreover/controller"
	"Moreover/middleware/auth"
)

func messageRouter() {
	r := Router.Group("/message").Use(auth.Auth())
	{
		r.GET("/connect", controller.HandleSSE)

		r.GET("/:action/:current/:pageSize", controller.GetMessages)
	}
}
