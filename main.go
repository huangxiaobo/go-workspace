package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/huangxiaobo/gospider/core/config"
	"github.com/huangxiaobo/gospider/core/log"
	"github.com/huangxiaobo/gospider/core/spider"
)

var ConfFilePath = flag.String("conf", "./etc/gospider.yml", "config file path")

func init() {
	// Init log
	log.InitLog("./output/", "gospider", "utf-8")
	// Parse command line arguments
	flag.Parse()
	// Load config
	config.LoadConfig(*ConfFilePath)
}

type DefaultParser struct{}

func (p *DefaultParser) Name() string {
	return "DefaultParser"
}

func (p *DefaultParser) Parse(content string) error {
	log.Info("parser content:", len(content))
	return nil
}

func main() {

	log.Info("GoSpider")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	// 恢复默认的中断信号行为, 然后通知spider退出
	defer stop()

	s := spider.NewSpider(ctx)

	s.AddUrl(
		"https://www.baidu.com/",
		&DefaultParser{},
	)

	// 启动spider
	go s.Run()

	// 监听中断退出信号.
	<-ctx.Done()

	if err := s.GracefullyShutdown(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

}
