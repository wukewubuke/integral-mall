package config

import (
	rpc_config "github.com/yakaa/grpcx/config"
)


type Config struct {

	RpcServerConfig *rpc_config.ServiceConf

	Mode string


	Mysql struct {
		DataSource string
		Table      struct {
			User string
		}
	}
	Redis struct {
		DataSource string
		Auth       string
	}
}
