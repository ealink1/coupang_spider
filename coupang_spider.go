// main entry for coupang_spider
package main

import (
	"flag"

	"coupang_spider/internal/config"
	"coupang_spider/internal/handler"
	"coupang_spider/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "./etc/coupang_spider.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	logx.Infof("Starting coupang_spider at %s:%d...", c.Host, c.Port)
	server.Start()
}
