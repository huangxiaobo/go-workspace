package main

import (
	"fmt"
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
				log.Info(fmt.Sprintf("download %s, status: %t ", reqUrl, ok))
				if ok != true {
					// retry
					fetchTasks <- fetchTask
					continue
				}
				log.Info("html: ", html)

				parser := parser.Factory(fetchTask.Parser)

				log.Info(fmt.Sprintf("crawler %s finish, parser:%s", reqUrl, parser))
				time.Sleep(time.Second)
			}

		}()

	}

}

func (sd *Scheduler) start() {
	fetchTasks := make(chan task.FetchTask, 100)

	go sd.crawler(fetchTasks)

	fetchTasks <- task.FetchTask{Url: "https://www.zhihu.com/people/huang-liao-57", Parser: "zhihu"}
	fetchTasks <- task.FetchTask{Url: "https://www.zhihu.com/people/gong-qing-tuan-zhong-yang-67", Parser: "zhihu"}
	fetchTasks <- task.FetchTask{Url: "https://www.zhihu.com/people/cloudycity", Parser: "zhihu"}
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
