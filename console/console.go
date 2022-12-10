package main

import (
	"flag"
	"fmt"

	"console-api/console/internal/config"
	"console-api/console/internal/handler"
	"console-api/console/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/console-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv()) // conf.UseEnv() 从环境变量中读取配置

	server := rest.MustNewServer(c.RestConf, rest.WithCors(c.Cors.AllowOrigin)) // rest.WithCors(c.CorsConf) 跨域配置

	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
