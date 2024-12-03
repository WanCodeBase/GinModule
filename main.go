package main

import (
	"github.com/gin-gonic/gin"
)

// 基础
func main() {
	ginServer := gin.Default()
	// ginServer.Use()

	ginServer.GET("/hello", func(context *gin.Context) {
		context.JSON(200, gin.H{"msg": "Hello World!"})
	})

	// RESTful API
	ginServer.GET("/user")
	ginServer.POST("/user")
	ginServer.PUT("/user")
	ginServer.DELETE("/user")

	ginServer.Run(":8082")
}
