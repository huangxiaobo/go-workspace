package main

import (
	"money/core/task"
	"money/core/utils"
	"time"

	"money/core/fetch"
	"money/core/log"
	"money/core/parser"
)

type Scheduler struct {
}

func (sd *Scheduler) crawler(fetchTasks chan task.FetchTask) {

	for i := 0; i < 10; i++ {
		go func() {
			for true {
				fetchTask := <-fetchTasks

				reqUrl := fetchTask.Url
				log.Info("crawler proxy: ", reqUrl)

				proxy := utils.GetProxy()

				ok, html := fetch.Fetch(reqUrl, proxy)
				log.Info("download ", reqUrl, " status: ", ok)
				if ok != true {
					// retry
					fetchTasks <- fetchTask
					continue
				}
				log.Info("html: ", html)

				parser := parser.Factory(fetchTask.Parser)

				log.Info("crawler %s finish, parser:%s", reqUrl, parser)
				time.Sleep(time.Second)
			}

		}()

	}

}

func (sd *Scheduler) start() {
	fetchTasks := make(chan task.FetchTask, 100)

	go sd.crawler(fetchTasks)

	fetchTasks <- task.FetchTask{Url: "https://rarbgmirror.org/torrent/q8c1nbv", Parser: "rarbg"}
}

func main() {
	log.InitLog("./output/money.log", "money", "utf-8")
	log.Info("proxy pool")

	sd := Scheduler{}
	sd.start()

	for {
		time.Sleep(time.Microsecond)
	}
}
