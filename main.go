package main

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"log"
	"money/core/fetch"
	"money/core/parser"
	"money/core/pool"
	"money/core/utils"
	"time"
)

type Scheduler struct {
	Client *redis.Client
}

func (sd *Scheduler) addFetchTasks(fetchTasks chan fetch.Task) {
	go func() {

		for i := 1; i < 12; i++ {
			reqUrl := fmt.Sprintf("https://www.kuaidaili.com/free/inha/%d/", i)
			fetchTasks <- fetch.Task{Url: reqUrl, Domain: utils.KUAIDAILI}
		}
	}()

}

func (sd *Scheduler) crawlerProxy(proxyPool *pool.ProxyPool, fetchTasks chan fetch.Task) {

	for i := 0; i < 10; i++ {
		go func() {
			for true {
				fetchTask := <-fetchTasks

				reqUrl := fetchTask.Url
				domain := fetchTask.Domain
				log.Printf("crawler proxy: %s\n", reqUrl)

				proxy := proxyPool.GetProxy()

				ok, html := fetch.Fetch(reqUrl, &proxy)
				fmt.Println("download ", reqUrl, " status: ", ok)
				if ok != true {
					// retry
					fetchTasks <- fetchTask
					continue
				}

				parser := parser.Factory(domain)
				ipList := parser.Parse(html)

				for i := 0; i < len(ipList); i++ {
					proxyPool.ReceiveProxy(ipList[i])
				}
				fmt.Printf("crawler %s finish\n", reqUrl)
				time.Sleep(time.Second)
			}

		}()

	}

}

func (sd *Scheduler) start() {
	var proxyPool = pool.ProxyPool{}
	proxyPool.Initialize()

	fetchTasks := make(chan fetch.Task, 100)

	sd.Client, _ = utils.Connect()

	go sd.addFetchTasks(fetchTasks)
	go sd.crawlerProxy(&proxyPool, fetchTasks)

}

func main() {
	fmt.Println("proxy pool")

	sd := Scheduler{}
	sd.start()

	for {
		time.Sleep(time.Microsecond)
	}
}
