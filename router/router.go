package router

import (
	"Moreover/setting"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"io"
	"time"
)

var Router *gin.Engine

func InitRouter(port string) {
	gin.SetMode(setting.Config.Mode)
	logPath := setting.Config.LogConfig.Path
	writer, _ := rotatelogs.New(
		logPath+"%Y-%m-%d.log",
		rotatelogs.WithLinkName(logPath),
		rotatelogs.WithMaxAge(time.Duration(180)*time.Second),
		rotatelogs.WithRotationTime(time.Duration(60)*time.Minute*24),
	)

	gin.DefaultWriter = io.MultiWriter(writer)

	Router = gin.Default()

	captchaRouter()

	userRouter()

	activityRouter()

	commentRouter()

	likeRouter()

	followRouter()

	panic(Router.Run(port))
}
