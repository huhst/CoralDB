package webBackstage

import (
	"CoralDB/server"
	"github.com/gin-gonic/gin"
)

var Server *server.Server

func InitRouter(router *gin.Engine, server *server.Server) *gin.Engine {
	Server = server
	// 主页
	router.GET("/admin", Admin)
	router.POST("/admin", AdminPost)

	router.GET("/user", User)
	router.POST("/user", UserPost)

	router.GET("/data", Data)
	router.POST("/data", DataPost)

	router.GET("/table", Table)
	router.POST("/table", TablePost)

	return router
}
