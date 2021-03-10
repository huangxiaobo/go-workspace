package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v7"
	log "github.com/sirupsen/logrus"
	"money/core/fetch"
	"money/core/parser"
	"money/core/pool"
	"money/core/utils"
	"time"
)

type Scheduler struct {
	redisClient *redis.Client
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
				log.Info("crawler proxy: %s\n", reqUrl)

				proxy := proxyPool.GetProxy()

				ok, html := fetch.Fetch(reqUrl, &proxy)
				log.Info("download ", reqUrl, " status: ", ok)
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

func (sd *Scheduler) validateProxy() {
	go func() {
		for true {
			// 从redis的队列中读取proxy 然后验证proxy是否有效
			var values []string
			values, err := sd.redisClient.BLPop(time.Duration(time.Duration.Seconds(5)), utils.REDIS_RAW_PROXY_LIST).Result()
			if err != nil {
				log.Println(err)
				continue
			}

			for _, value := range values {
				log.WithFields(log.Fields{"proxy": value})
				var proxy utils.Proxy
				err = json.Unmarshal([]byte(value), &proxy)
				if err != nil {
					log.Println(err)
					continue
				}
				log.Printf("hxb >>>>>>>>>%v\n", proxy.Ip)

			}

			time.Sleep(time.Duration(time.Duration.Seconds(2)))
		}
	}()
}

func (sd *Scheduler) start() {
	var proxyPool = pool.ProxyPool{}
	proxyPool.Initialize()

	fetchTasks := make(chan fetch.Task, 100)

	sd.redisClient, _ = utils.Connect()

	go sd.addFetchTasks(fetchTasks)
	go sd.crawlerProxy(&proxyPool, fetchTasks)
	go sd.validateProxy()

}

func main() {
	utils.InitLog("./log", "money", "utf-8")
	log.Info("proxy pool")

	sd := Scheduler{}
	sd.start()

	for {
		time.Sleep(time.Microsecond)
	}
}
