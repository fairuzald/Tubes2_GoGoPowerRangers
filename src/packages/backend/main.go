package main

import (
	"backend/handlers"
	"backend/router"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	handlers.InitCache()
	r := router.SetupRouter()

	r.Use(CORSMiddleware())

	r.Run()
	defer func() {
		if err := handlers.SaveCacheToJSON("cache.json"); err != nil {
			fmt.Println("Error saving cache to JSON:", err)
		}
	}()
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
