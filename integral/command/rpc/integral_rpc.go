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
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"my-integral-mall/integral/command/rpc/config"
	"my-integral-mall/integral/logic"
	"my-integral-mall/integral/model"
	"my-integral-mall/integral/protos"
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

	integralModel := model.NewIntegralModel(engine, client, conf.Mysql.Table.Integral)
	integralLogic, err := logic.NewIntegralLogic(conf.RabbitMq.DataSource+conf.RabbitMq.VirtualHost,
		conf.RabbitMq.QueueName, integralModel)
	if err != nil {
		panic(err)
	}
	defer integralLogic.CloseRabbitMqConn()

	rpcServer, err := grpcx.MustNewGrpcxServer(conf.RpcServerConfig, func(server *grpc.Server) {
		protos.RegisterIntegralRpcServer(server, integralLogic)
	})

	if err != nil {
		log.Fatal(err)
	}

	integralLogic.ConsumeMessage()
	//integralLogic.PushMessage("insert into integral(user_id,integral)values(10,20)")

	log4g.Error(rpcServer.Run())

}
