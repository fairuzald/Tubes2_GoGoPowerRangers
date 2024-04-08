package router

import (
	"backend/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", handlers.PingHandler)
	r.POST("/algorithm", handlers.PostAlgorithmHandler)
	r.POST("/scraping", handlers.ScrapingHandlerPost)
	r.GET("/autocomplete", handlers.AutoCompleteHandler)
	return r
}
