package main

import (
	"common/lib"
	"config/configHandle"
	"flag"
	"fmt"
	"micro_tmpl/player/config"
	"micro_tmpl/player/handler"
)

var port string
var env string

func main() {
	flag.StringVar(&port, "p", "50051", "端口号")
	flag.StringVar(&env, "env", "local", "环境")
	flag.Parse()

	//读取配置
	//全局配置
	_, err := configHandle.LoadAllConfig(env)
	if err != nil {
		return
	}
	//单独配置
	_, err = config.LoadAllConfig()
	if err != nil {
		return
	}

	//连接mongodb
	mongoConfig := configHandle.GetConfigManInstance().GetGeneralConfig().MongodbConfig
	mongo := lib.GetMongoStore()
	err = mongo.CreateConn(mongoConfig.ConnAddress, mongoConfig.DBName)
	if err != nil {
		return
	}

	rpcServer := handler.ProtoServer{
		Mongo: mongo,
	}
	address := fmt.Sprintf("127.0.0.1:%s", port)
	go rpcServer.RegisterRpcServer(address)
	go rpcServer.RegisterEtcdServer(address)
	fmt.Printf("listen address %s\n", address)
	select {}
}
