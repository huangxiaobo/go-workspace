package pool

import (
	"../utils"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type ProxyPool struct {
	out chan utils.Proxy
	in  chan utils.Proxy
}

const ProxyCsvFile = "proxy_pool.csv"

func (pp *ProxyPool) Initialize() *ProxyPool {
	pp.out = make(chan utils.Proxy)
	pp.in = make(chan utils.Proxy)

	pp.readOut()
	pp.writeIn()

	return pp
}

func (pp *ProxyPool) readOut() {
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

	pp.in <- proxy

}
