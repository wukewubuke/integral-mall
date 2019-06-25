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
	 "my-integral-mall/common/middleware"
	"my-integral-mall/common/rpcxclient/integralrpcmodel"
	"my-integral-mall/goods/command/api/config"
	"my-integral-mall/goods/controller"
	"my-integral-mall/goods/logic"
	"my-integral-mall/goods/model"
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



	goodsModel := model.NewGoodsModel(engine, client, conf.Mysql.Table.Goods)
	goodsLogic := logic.NewGoodsLogic(goodsModel, integralRpcModel)
	goodsController := controller.NewGoodsController(goodsLogic)

	middleware := middleware.NewAuthorization(client)


	r := gin.Default()
	r.Use(middleware.Auth)
	goodsGroup := r.Group("/goods")
	{
		goodsGroup.POST("/search",goodsController.GoodsSearch)
		goodsGroup.POST("/list",goodsController.GoodsSearch)
		goodsGroup.POST("/order",goodsController.GoodsOrder)
	}



	log4g.Error(r.Run(conf.Port))
}
