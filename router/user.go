package router

import (
	"Moreover/controller"
	"Moreover/middleware/auth"
)

func userRouter() {
	r := Router.Group("/user")
	{
		r.POST("/register", controller.Register)

		r.POST("/login", controller.Login)

		r.GET("/info/:userId", auth.Auth(), controller.GetUserInfoById)

		r.PATCH("/:info", auth.Auth(), controller.UpdateUserInfo)
	}
}
