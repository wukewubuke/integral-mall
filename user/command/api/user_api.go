package main

import (
	"github.com/gin-gonic/gin"
	"github.com/yakaa/log4g"
)

func main() {

	log4g.Init(log4g.Config{Path:"logs"})
	gin.DefaultWriter = log4g.InfoLog
	gin.DefaultErrorWriter = log4g.ErrorLog
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	log4g.Error(r.Run(":8888")) // listen and serve on 0.0.0.0:8080
}