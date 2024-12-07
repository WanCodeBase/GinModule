package main

import (
	"database/sql"
	"github.com/WanCodeBase/GinModule/util"
	"log"

	"github.com/WanCodeBase/GinModule/api"
	db "github.com/WanCodeBase/GinModule/db/sqlc"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// 基础
func _main() {
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

func main() {
	conf, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalln("load config failed:", err)
		return
	}
	conn, err := sql.Open(conf.DBDriver, conf.DBSource)
	if err != nil {
		log.Fatalln("connect db failed:", err)
		return
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(conf.ServerAddress)
	if err != nil {
		log.Fatalln("server start failed:", err)
		return
	}
}
