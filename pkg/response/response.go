package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Response(ctx *gin.Context, code int, data gin.H) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": GetMessage(code),
		"data":    data,
	})
}
