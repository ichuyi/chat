package handler

import (
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	//gin.SetMode(gin.ReleaseMode)
	r.POST("/login", login)
	r.Use(verify)
	r.Any("/chat", chat)
	return r
}
