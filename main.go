package main

import (
	"flag"

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
	log.Info("parser content:", content)
	return nil
}

func main() {

	log.Info("Crawler")

	sd := spider.NewSpider()

	sd.AddFetchTask(&spider.FetchTask{
		Url:    "https://www.baidu.com/",
		Parser: &DefaultParser{},
	})

	sd.Start()

}
