package fetch

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	utils2 "money/pkg/crawler/core/utils"
	"net/http"
	"net/url"
	"time"
)

type Task struct {
	Url     string
	Domain  utils2.DomainType
	Content string
}

/**
* 随机返回一个User-Agent
 */
func getAgent() string {
	agent := [...]string{
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

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	len := len(agent)
	return agent[r.Intn(len)]
}

func Fetch(urlString string, proxy *utils2.Proxy) (bool, string) {
	fmt.Printf("download >>> url: %s, proxy: %v\n", urlString, *proxy)

	timeout := time.Duration(5000 * time.Millisecond) //超时时间50ms

	proxyString := fmt.Sprintf("%s://%s:%s", proxy.Protocol, proxy.Ip, proxy.Port)
	proxyUrl, _ := url.Parse(proxyString)

	transport := &http.Transport{
		Proxy:           http.ProxyURL(proxyUrl),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: transport, Timeout: timeout}
	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		log.Fatal("new request failed,", err.Error())
		return false, ""
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", getAgent())
	//req.Header.Set("Connection", "keep-alive")

	resp, err := client.Do(req)

	if err != nil || resp == nil {
		fmt.Printf("do request fail>>>: %v", err)
		return false, ""
	}

	fmt.Println(resp.StatusCode)
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return false, ""
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("read body failed:" + err.Error())
	}

	return true, string(body)
}
