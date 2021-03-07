package pool

import (
	"bufio"
	"fmt"
	"github.com/go-redis/redis/v7"
	"log"
	"money/core/utils"
	"os"
	"strings"
	"time"
)

type ProxyPool struct {
	out chan utils.Proxy
	in  chan utils.Proxy
	Client *redis.Client
}

const ProxyCsvFile = "proxy_pool.csv"

func (pp *ProxyPool) Initialize() *ProxyPool {
	pp.out = make(chan utils.Proxy)
	pp.in = make(chan utils.Proxy)

	pp.readOut()
	pp.writeIn()

	pp.Client, _ = utils.Connect()


	return pp
}

func (pp *ProxyPool) load()  {
	// 将cvs文件的内容加载到redis
	file, err := os.OpenFile(ProxyCsvFile, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	reader := bufio.NewReader(file)

	for true {
		line, err := reader.ReadString('\n')
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		line = utils.Strip(line)
		arr := strings.Split(line, ",")

		if len(arr) < 5 {
			fmt.Printf("line is not valid: %s\n", line)
			continue
		}

		proxy := utils.Proxy{Ip: arr[0], Port: arr[1], Anonymous: arr[2], Protocol: arr[3], Location: arr[4]}
		pp.Client.LPush("proxy", proxy)
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
			if err != nil {
				log.Fatal(err)
				time.Sleep(5 * time.Second)
				continue
			}

			line = utils.Strip(line)
			fmt.Printf(" > Read line %s \n", line)

			arr := strings.Split(line, ",")

			if len(arr) < 5 {
				fmt.Printf("line is not valid: %s\n", line)
				continue
			}

			proxy := utils.Proxy{Ip: arr[0], Port: arr[1], Anonymous: arr[2], Protocol: arr[3], Location: arr[4]}
			log.Printf(proxy.String())
			pp.out <- proxy
		}

	}()
}

func (pp *ProxyPool) GetProxy() utils.Proxy {

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

func (pp *ProxyPool) ReceiveProxy(proxy utils.Proxy) {
	// 收到解析的代理

	pp.in <- proxy
	pp.Client.LPush("proxy", proxy)

}
