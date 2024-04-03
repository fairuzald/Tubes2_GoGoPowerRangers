package router

import (
	"backend/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", handlers.PingHandler)

	r.POST("/api/algorithm", handlers.PostAlgorithmHandler)

	return r
}
