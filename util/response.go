package util

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type CommonResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func OKResponse(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, CommonResponse{
		Code:    0,
		Message: "ok",
		Data:    data,
	})
}
func FailedResponse(ctx *gin.Context, code int, message string) {
	ctx.JSON(http.StatusOK, CommonResponse{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}
