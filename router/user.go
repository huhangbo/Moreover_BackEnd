package router

import (
	"Moreover/controller"
	"Moreover/middleware"
)

func userRouter() {
	r := Router.Group("/user")
	{
		r.POST("/register", controller.Register)

		r.POST("/login", controller.Login)

		r.GET("/info/:userId", middleware.Auth(), controller.GetUserInfoById)

		r.PATCH("/:info", middleware.Auth(), controller.UpdateUserInfo)
	}
}
