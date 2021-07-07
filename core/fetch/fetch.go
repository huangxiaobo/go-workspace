package fetch

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"money/core/db"
	"money/core/log"
	"money/core/utils"
)

var agents = []string{
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:50.0) Gecko/20100101 Firefox/50.0",
	"Opera/9.80 (Macintosh; Intel Mac OS X 10.6.8; U; en) Presto/2.8.131 Version/11.11",
	"Opera/9.80 (Windows NT 6.1; U; en) Presto/2.8.131 Version/11.11",
	"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; 360SE)",
	"Mozilla/5.0 (Windows NT 6.1; rv:2.0.1) Gecko/20100101 Firefox/4.0.1",
	"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; The World)",
	"User-Agent,Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_6_8; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
	"User-Agent, Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Maxthon 2.0)",
	"User-Agent,Mozilla/5.0 (Windows; U; Windows NT 6.1; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
}

/**
* 随机返回一个User-Agent
 */
func getAgent() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	size := len(agents)
	return agents[r.Intn(size)]
}

func Fetch(urlString string, proxy *utils.ProxyObj) (bool, string) {
	log.Info(fmt.Sprintf("download >>> url: %s, proxy: %+v", urlString, proxy))

	transport := proxy.GetTransport()
	client := &http.Client{Transport: transport, Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		log.Fatal("new request failed,", err.Error())
		return false, ""
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", getAgent())
	// req.Header.Set("Connection", "keep-alive")

	resp, err := client.Do(req)

	if err != nil || resp == nil {
		log.InfoWithFields("do request fail>>>: ", log.Fields{"err": err})
		return false, ""
	}

	log.InfoWithFields("", log.Fields{"url": urlString, "status_code": resp.StatusCode})
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return false, ""
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("read body failed:", err.Error())
	}

	return true, string(body)
}

type Fetcher struct {
}

func (f *Fetcher) Start() {
	for true {
		task := &db.FetchTask{}
		if err := db.GetFetchTask(task); err != nil {
			time.Sleep(10 * time.Second)
			continue
		}

		reqUrl := task.Url
		log.Info("crawler proxy: ", reqUrl)

		proxy := utils.GetProxy()

		ok, html := Fetch(reqUrl, proxy)
		log.Info(fmt.Sprintf("download %s, status: %t ", reqUrl, ok))
		if ok != true {
			// retry
			continue
		}
		log.Info("html: ", html)

		task.Page = html

		if err := db.UpdateFetchTask(task); err != nil {
			logrus.Error(err)
		}

		time.Sleep(time.Second)
	}
}
