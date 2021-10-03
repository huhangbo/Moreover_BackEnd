package auth

import (
	"Moreover/pkg/jwt"
	"Moreover/pkg/response"
	"github.com/gin-gonic/gin"
	"strings"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			response.Response(c, response.AuthError, nil)
			c.Abort()
			return
		}
		claim := jwt.ParseToken(strings.TrimPrefix(tokenString, "Bearer "))
		if claim == nil {
			response.Response(c, response.AuthError, nil)
			c.Abort()
			return
		}
		StuId := claim.StuId
		c.Set("stuId", StuId)
		return
	}
}

func Verify() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userId")
		stuId, ok := c.Get("stuId")
		if ok != true || userId != stuId.(string) {
			response.Response(c, response.AuthError, nil)
			c.Abort()
		}
	}
}
