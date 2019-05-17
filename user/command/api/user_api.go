package main

import (
	"encoding/json"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/yakaa/log4g"
	"io/ioutil"
	"log"
	"my-integral-mall/user/command/api/config"
)



var configFile = flag.String("f","config/config.json","use config")


func main() {
	flag.Parse()
	var conf config.Config
	bs, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(bs, &conf); err != nil {
		log.Fatal(err)
	}


	if conf.Mode == gin.ReleaseMode {
		log4g.Init(log4g.Config{Path:"logs"})
		gin.DefaultWriter = log4g.InfoLog
		gin.DefaultErrorWriter = log4g.ErrorLog
	}


	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	//gr.Run()
	log4g.Error(r.Run(conf.Port))
}