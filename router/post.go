package router

import (
	"Moreover/controller"
	"Moreover/middleware"
)

func PostRouter() {
	r := Router.Group("/post", middleware.Auth())
	{
		r.POST("/", controller.PublishPost)

		r.GET("/:type/:current/:pageSize", controller.GetPostByPage)

		r.PUT("/:postId", controller.UpdatePost)

		r.DELETE("/:postId", controller.DeletePost)

		r.GET("/follow/:current/:pageSize", controller.GetFollowPostByPage)

	}
}
