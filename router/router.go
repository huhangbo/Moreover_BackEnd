package router

import (
	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

func InitRouter(port string)  {
	Router = gin.Default()
	UserRouter()
	panic(Router.Run(port))
}