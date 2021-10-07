package router

import (
	"Moreover/controller/captchaController"
)

func captchaRouter() {
	r := Router.Group("/captcha")
	{
		r.GET("/generate", captchaController.GenerateCaptcha)

		r.POST("/parse", captchaController.ParseCaptcha)
	}
}
