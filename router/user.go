package router

import (
	"Moreover/controller/userController"
	"Moreover/middleware/auth"
)

func UserRouter() {
	r := Router.Group("/user")
	{
		r.POST("/register", userController.Register)

		r.POST("/login", userController.Login)

		r.GET("/info/:userId", userController.GetUserInfoById)

		r.PATCH("/avatar", auth.Auth(), userController.UpdateAvatar)

		r.PATCH("/sex", auth.Auth(), userController.UpdateSex)

		r.PATCH("/nickname", auth.Auth(), userController.UpdateNickname)

		r.PATCH("/description", auth.Auth(), userController.UpdateDescription)
	}

}
