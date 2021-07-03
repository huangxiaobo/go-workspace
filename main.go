package main

import (
	"flag"
	"time"

	"money/core/config"
	"money/core/db"
	"money/core/fetch"
	"money/core/log"
)

type Scheduler struct {
}

func (sd *Scheduler) crawler() {

	go func() {
		fetcher := &fetch.Fetcher{}
		fetcher.Start()
	}()

}

func (sd *Scheduler) start() {

	go sd.crawler()

	for _, item := range config.Conf.Crawler.Tasks {
		fetchTask := &db.FetchTask{
			Project: item.Project,
			Url:     item.Url,
			Parser:  item.Parser,
		}
		db.AddFetchTask(fetchTask)

	}
}

var ConfFilePath = flag.String("conf", "./config/money.yml", "config file path")

func init() {
	// Init log
	log.InitLog("./output/", "money", "utf-8")
	// Parse command line arguments
	flag.Parse()
	// Load config
	config.LoadConfig(*ConfFilePath)
	// Database SetUp
	db.SetUp(&config.Conf)
}

func main() {

	log.Info("Crawler")

	sd := Scheduler{}
	sd.start()

	for {
		time.Sleep(time.Microsecond)
	}

}
