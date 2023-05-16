package main

import (
	"common/lib"
	"config/configHandle"
	"flag"
	"fmt"
	"micro_tmpl/gateway/route"
)

func main() {
	env := "local"
	port := "8001"
	//解析命令行参数
	flag.StringVar(&env, "env", "local", "执行环境")
	flag.StringVar(&port, "port", "8001", "端口")
	flag.Parse()
	//加载配置
	_, err := configHandle.LoadAllConfig(env)
	if err != nil {
		return
	}

	//初始化日志:
	lib.InitLogger()

	//启动服务器
	address := fmt.Sprintf("127.0.0.1:%s", port)
	appSecret := configHandle.GetConfigManInstance().GetGeneralConfig().AppSecret
	apiServer := lib.CreateHttpServer(address, appSecret)
	apiServer.Run()
	apiRoute := route.GetRouteInstance()
	apiServer.RegisterRoutes("/:module/:method", apiRoute.DispatchApiRequest)
	lib.Logger.Infof("AppStart: ApiServer Listen %s", address)
	select {}
}
