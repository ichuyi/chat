package handler

import (
	"chat/util"
	"github.com/gin-gonic/gin"
)

func verify(ctx *gin.Context) {
	if username, err := ctx.Cookie("username"); err != nil {
		util.FailedResponse(ctx, util.Unauthorized, util.UnauthorizedMsg)
		ctx.Abort()
	} else {
		ctx.Set("name", username)
		ctx.Next()
	}
}
