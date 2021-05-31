package parser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	utils2 "money/pkg/crawler/core/utils"
	"strings"
)

type Kuaidaili struct {
}

func (parser *Kuaidaili) Parse(content string) []utils2.Proxy {

	var proxyList []utils2.Proxy

	dom, err := goquery.NewDocumentFromReader(strings.NewReader(content))

	if err != nil {
		log.Fatal("goquery read failed:" + err.Error())
	}

	table := dom.Find("div#list").Find("table.table.table-bordered").Eq(0)

	tableHead := table.Find("thead").Eq(0)
	tableHead.Find("tr").Find("th").Each(func(i int, selection *goquery.Selection) {
		fmt.Println(">>>", i)
		fmt.Println(utils2.Strip(selection.Text()))
		fmt.Println("<<<", i)

	})

	tableBody := table.Find("tbody").Eq(0)
	tableBody.Find("tr").Each(func(i int, selection *goquery.Selection) {

		tds := selection.Find("td")

		proxyIp := utils2.Strip(tds.Eq(0).Text())
		fmt.Println(proxyIp)

		proxyPort := utils2.Strip(tds.Eq(1).Text())
		fmt.Println(proxyPort)

		anonymousType := utils2.Strip(tds.Eq(2).Text())
		fmt.Println("匿名:", anonymousType)

		proxyType := utils2.Strip(tds.Eq(3).Text())
		fmt.Println("类型:", proxyType)

		proxyLocation := utils2.Strip(tds.Eq(4).Text())
		fmt.Println("位置", proxyLocation)

		proxySpeed := utils2.Strip(tds.Eq(5).Text())
		fmt.Println("响应速度:", proxySpeed)

		proxyLastCheck := utils2.Strip(tds.Eq(6).Text())
		fmt.Println("最后验证时间::", proxyLastCheck)

		proxy := utils2.Proxy{Ip: proxyIp, Port: proxyPort, Anonymous: anonymousType, Protocol: proxyType, Location: proxyLocation}

		proxyList = append(proxyList, proxy)
	})

	return proxyList
}
