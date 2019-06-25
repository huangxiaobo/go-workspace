package src

import (
	"./fetch"
	"./parser"
	"./pool"
	"./utils"
	"fmt"
	"time"
)

func addFetchTasks(fetchTasks chan fetch.Task) {
	go func() {

		for i := 1; i < 12; i++ {
			reqUrl := fmt.Sprintf("https://www.kuaidaili.com/free/inha/%d/", i)
			fetchTasks <- fetch.Task{Url: reqUrl, Domain: utils.KUAIDAILI}
		}
	}()

	go func() {
		go func() {
			for i := 1; i < 15; i++ {
				reqUrl := fmt.Sprintf("https://www.xicidaili.com/nn/%d", i)
				fetchTasks <- fetch.Task{Url: reqUrl, Domain: utils.XICI}
			}
		}()

		go func() {
			for i := 1; i < 10; i++ {
				reqUrl := fmt.Sprintf("https://www.xicidaili.com/nt/%d", i)
				fetchTasks <- fetch.Task{Url: reqUrl, Domain: utils.XICI}
			}
		}()

		go func() {
			for i := 1; i < 10; i++ {
				reqUrl := fmt.Sprintf("https://www.xicidaili.com/wn/%d", i)
				fetchTasks <- fetch.Task{Url: reqUrl, Domain: utils.XICI}
			}
		}()

		go func() {
			for i := 1; i < 10; i++ {
				reqUrl := fmt.Sprintf("https://www.xicidaili.com/wt/%d", i)
				fetchTasks <- fetch.Task{Url: reqUrl, Domain: utils.XICI}
			}
		}()

	}()

}

func crawlerProxy(proxyPool *pool.ProxyPool, fetchTasks chan fetch.Task) {

	for i := 0; i < 10; i++ {
		go func() {
			for true {
				fetchTask := <-fetchTasks

				reqUrl := fetchTask.Url
				domain := fetchTask.Domain

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

func main() {
	fmt.Println("proxy pool")

	var proxyPool = pool.ProxyPool{}
	proxyPool.Initialize()

	fetchTasks := make(chan fetch.Task, 100)

	go addFetchTasks(fetchTasks)
	go crawlerProxy(&proxyPool, fetchTasks)

	for {
		time.Sleep(time.Microsecond)
	}
}
