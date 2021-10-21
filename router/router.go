package router

import (
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

	captchaRouter()

	userRouter()

	activityRouter()

	commentRouter()

	likeRouter()

	followRouter()

	PostRouter()

	panic(Router.Run(port))
}
