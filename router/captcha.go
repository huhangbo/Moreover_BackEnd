package router

import "Moreover/controller"

func CaptchaRouter() {
	r := Router.Group("/captcha")
	{
		r.GET("/generate", controller.GenerateCaptcha)

		r.POST("/parse", controller.ParseCaptcha)
	}
}
