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

		r.PATCH("/avatar/:userId", auth.Auth(), auth.Verify(), userController.UpdateAvatar)

		r.PATCH("/sex/:userId", auth.Auth(), auth.Verify(), userController.UpdateSex)

		r.PATCH("/nickname/:userId", auth.Auth(), auth.Verify(), userController.UpdateNickname)

		r.PATCH("/description/:userId", auth.Auth(), auth.Verify(), userController.UpdateDescription)
	}

}
