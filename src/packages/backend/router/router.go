package router

import (
	"backend/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", handlers.PingHandler)
	r.POST("/url-info", handlers.URLInfoHandler)
	r.GET("/autocomplete", handlers.AutoCompleteHandler)

	r.POST("/bfs", handlers.BFSHTTPHandler)
	r.POST("/bfs/backup", handlers.BFSHTTPHandlerBackup)
	r.POST("/ids", handlers.IDSHTTPHandler)
	r.POST("/ids/backup", handlers.IDSHTTPHandlerBackup)
	r.POST("/save", handlers.SaveCacheHandler)
	return r
}
