package utils

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

//
const s = `
{
  "anonymous": "",
  "check_count": 3,
  "fail_count": 0,
  "https": false,
  "last_status": true,
  "last_time": "2021-06-14 17:22:42",
  "proxy": "106.52.10.171:9999",
  "region": "",
  "source": "freeProxy09"
}
`

type ProxyObj struct {
	Proxy      string `json:"proxy"`
	Anonymous  string `json:"anonymous"`
	Https      bool   `json:"https"`
	LastStatus bool   `json:"last_status"`
}

func (po *ProxyObj) GetProtocol() string {
	if po.Https {
		return "https"
	} else {
		return "http"
	}
}

func (po *ProxyObj) GetTransport() *http.Transport {
	proxyString := fmt.Sprintf("%s://%s", po.GetProtocol(), po.Proxy)
	proxyUrl, _ := url.Parse(proxyString)

	transport := &http.Transport{
		Proxy:           http.ProxyURL(proxyUrl),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return transport
}

func (po *ProxyObj) String() string {
	return fmt.Sprintf("Proxy{https=%b, port=%s}", po.Https, po.Proxy)
}

func GetProxy() *ProxyObj {
	url := "http://127.0.0.1:5010/get/"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("err : %s", err)
	}

	data, err := ioutil.ReadAll(resp.Body)

	proxyObj := &ProxyObj{}
	json.Unmarshal(data, &proxyObj)
	return proxyObj

}
