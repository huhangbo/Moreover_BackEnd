package router

import (
	"Moreover/middleware"
	"Moreover/setting"
	"github.com/gin-gonic/gin"
	"io"
	"os"
)

var Router *gin.Engine

func InitRouter(port string) {

	gin.SetMode(setting.Config.Mode)

	f, _ := os.Create("log/gin.log")

	gin.DefaultWriter = io.MultiWriter(f)

	Router = gin.Default()

	Router.Use(middleware.Cors())

	captchaRouter()

	userRouter()

	activityRouter()

	commentRouter()

	likeRouter()

	followRouter()

	PostRouter()

	messageRouter()

	panic(Router.Run(port))
}
