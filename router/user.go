package router

import "Moreover/controller"

func UserRouter() {
	r := Router.Group("/user")
	{

		r.POST("/register", controller.Register)

		r.POST("/login", controller.Login)
	}

}
