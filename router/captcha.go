package router

import (
	"Moreover/controller/captchaController"
)

func CaptchaRouter() {
	r := Router.Group("/captcha")
	{
		r.GET("/generate", captchaController.GenerateCaptcha)

		r.POST("/parse", captchaController.ParseCaptcha)
	}
}
