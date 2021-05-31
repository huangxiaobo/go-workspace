package parser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	utils2 "money/pkg/crawler/core/utils"
	"strings"
)

type Xicidaili struct {
}

func (parser *Xicidaili) Parse(content string) []utils2.Proxy {
	var ipList []utils2.Proxy

	dom, err := goquery.NewDocumentFromReader(strings.NewReader(content))

	if err != nil {
		log.Fatal(err.Error())
	}

	dom.Find("table").Find("tr").Each(func(i int, selection *goquery.Selection) {
		tds := selection.Find("td")
		if tds.Size() < 10 {
			return
		}

		proxyIp := utils2.Strip(tds.Eq(1).Text())
		fmt.Println(proxyIp)

		proxyPort := utils2.Strip(tds.Eq(2).Text())
		fmt.Println(proxyPort)

		proxyLocation := utils2.Strip(tds.Eq(3).Text())
		fmt.Println("位置", proxyLocation)

		anonymousType := utils2.Strip(tds.Eq(4).Text())
		fmt.Println("匿名:", anonymousType)

		proxyProtocol := utils2.Strip(tds.Eq(5).Text())
		fmt.Println("类型:", proxyProtocol)

		proxySpeed := utils2.Strip(tds.Eq(8).Text())
		fmt.Println("响应速度:", proxySpeed)

		proxyLastCheck := utils2.Strip(tds.Eq(9).Text())
		fmt.Println("最后验证时间::", proxyLastCheck)

		proxy := utils2.Proxy{Ip: proxyIp, Port: proxyPort, Anonymous: anonymousType, Protocol: proxyProtocol, Location: proxyLocation}

		ipList = append(ipList, proxy)
	})

	return ipList
}
