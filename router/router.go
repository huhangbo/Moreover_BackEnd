package router

import (
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"io"
	"time"
)

var Router *gin.Engine

func InitRouter(port string) {
	path := "log/"
	writer, _ := rotatelogs.New(
		path+"%Y-%m-%d.log",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Duration(180)*time.Second),

		rotatelogs.WithRotationTime(time.Duration(60)*time.Minute*24),
	)

	gin.DefaultWriter = io.MultiWriter(writer)

	Router = gin.Default()

	//Router = gin.New()

	//Router.Use(logger.Logger()).Use(gin.Recovery())
	CaptchaRouter()

	UserRouter()

	ActivityRouter()

	panic(Router.Run(port))
}
