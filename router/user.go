package router

import (
	"Moreover/controller/userController"
)

func UserRouter() {
	r := Router.Group("/user")
	{

		r.POST("/register", userController.Register)

		r.POST("/login", userController.Login)
	}

}
