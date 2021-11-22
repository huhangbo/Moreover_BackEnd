package router

import (
	"Moreover/controller"
	"Moreover/middleware/auth"
)

func PostRouter() {
	r := Router.Group("/post", auth.Auth())
	{
		r.POST("/", controller.PublishPost)

		r.GET("/:type/:current/:pageSize", controller.GetPostByPage)

		r.PUT("/:postId", controller.UpdatePost)

		r.DELETE("/:postId", controller.DeletePost)

		r.GET("/follow/:current/:pageSize", controller.GetFollowPostByPage)

	}
}
