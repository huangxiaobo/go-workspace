package pool

import (
	"bufio"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	utils2 "money/pkg/crawler/core/utils"
	"os"
	"strings"
)

type ProxyPool struct {
	out    chan utils2.Proxy
	in     chan utils2.Proxy
	Client *redis.Client
}

const ProxyCsvFile = "proxy_pool.csv"

func (pp *ProxyPool) Initialize() *ProxyPool {
	pp.out = make(chan utils2.Proxy)
	pp.in = make(chan utils2.Proxy)

	pp.readOut()
	pp.writeIn()

	pp.Client, _ = utils2.Connect()

	pp.load2Redis()

	return pp
}

func (pp *ProxyPool) load2Redis() {
	// 将cvs文件的内容加载到redis
	file, err := os.OpenFile(ProxyCsvFile, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	reader := bufio.NewReader(file)

	for true {
		line, err := reader.ReadString('\n')
		if len(line) == 0 && err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		line = utils2.Strip(line)
		arr := strings.Split(line, ",")

		if len(arr) < 5 {
			fmt.Printf("line is not valid: %s\n", line)
			continue
		}

		proxy := utils2.Proxy{Ip: arr[0], Port: arr[1], Anonymous: arr[2], Protocol: arr[3], Location: arr[4]}

		proxyJsonData, _ := json.Marshal(proxy)
		log.Printf(">>>>>%v\n", proxyJsonData)
		value, err := pp.Client.LPush(utils2.REDIS_RAW_PROXY_LIST, proxyJsonData).Result()
		log.Println(err)
		log.Printf("send %s to redis: %d\n", proxy.Ip, value)
	}

}

func (pp *ProxyPool) readOut() {
	log.Printf("readout")
	go func() {
		file, err := os.OpenFile(ProxyCsvFile, os.O_RDONLY|os.O_CREATE, os.ModePerm)
		if err != nil {
			log.Fatal("proxy record file is not exists.")
		}

		defer file.Close()

		reader := bufio.NewReader(file)

		for true {
			line, err := reader.ReadString('\n')
			if len(line) == 0 && err != nil {
				if err == io.EOF {
					break
				}
				log.Fatal(err)
			}

			line = utils2.Strip(line)
			fmt.Printf(" > Read line %s \n", line)

			arr := strings.Split(line, ",")

			if len(arr) < 5 {
				fmt.Printf("line is not valid: %s\n", line)
				continue
			}

			proxy := utils2.Proxy{Ip: arr[0], Port: arr[1], Anonymous: arr[2], Protocol: arr[3], Location: arr[4]}
			log.Printf(proxy.String())
			pp.out <- proxy
		}

	}()
}

func (pp *ProxyPool) GetProxy() utils2.Proxy {

	return <-pp.out

}

func (pp *ProxyPool) writeIn() {
	go func() {
		file, err := os.OpenFile(ProxyCsvFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModePerm)
		if err != nil {
			log.Fatal(err.Error())
		}
		defer file.Close()
		for true {
			proxy := <-pp.in

			line := fmt.Sprintf("%s,%s,%s,%s,%s\n", proxy.Ip, proxy.Port, proxy.Anonymous, proxy.Protocol, proxy.Location)
			if _, err := file.WriteString(line); err != nil {
				log.Println(err)
			}
		}
	}()

}

func (pp *ProxyPool) ReceiveProxy(proxy utils2.Proxy) {
	// 收到解析的代理

	pp.in <- proxy
	pp.Client.LPush("proxy", proxy)

}
