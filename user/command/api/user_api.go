package main

import (
	"encoding/json"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/yakaa/grpcx"
	"github.com/yakaa/log4g"
	"io/ioutil"
	"log"
	"my-integral-mall/common/rpcxclient/integralrpcmodel"
	"my-integral-mall/user/command/api/config"
	"my-integral-mall/user/controller"
	"my-integral-mall/user/logic"
	"my-integral-mall/user/model"
)

var configFile = flag.String("f", "config/config.json", "use config")

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

	log4g.Init(log4g.Config{Path: "logs"})
	gin.DefaultWriter = log4g.InfoLog
	gin.DefaultErrorWriter = log4g.ErrorLog

	//初始化mysql  xorm
	engine, err := xorm.NewEngine("mysql", conf.Mysql.DataSource)
	if err != nil {
		log.Fatal(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.DataSource,
		Password: conf.Redis.Auth,
	})

	rpcClient, err := grpcx.MustNewGrpcxClient(conf.IntegralRpc)
	if err != nil {
		log.Fatal(err)
	}

	integralRpcModel := integralrpcmodel.NewIntegralRpcModel(rpcClient)



	userModel := model.NewUserModel(engine, client, conf.Mysql.Table.User)
	userLogic := logic.NewUserLogic(userModel, client, integralRpcModel)
	userController := controller.NewUserController(userLogic)



	r := gin.Default()
	userGroup := r.Group("/user")
	{
		userGroup.POST("/register",userController.Register)
		userGroup.POST("/login",userController.Login)
	}



	log4g.Error(r.Run(conf.Port))
}
