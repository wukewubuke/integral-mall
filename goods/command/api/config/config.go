package config

import (
	rpc_config "github.com/yakaa/grpcx/config"
)

type Config struct {
	Mode string
	Port string

	Mysql struct {
		DataSource string
		Table      struct {
			Goods string
		}
	}

	Redis struct {
		DataSource string
		Auth       string
	}
	RabbitMq struct {
		DataSource  string
		VirtualHost string
	}
	IntegralRpc *rpc_config.ClientConf
}
