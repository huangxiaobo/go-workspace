package utils

import "fmt"

type Proxy struct {
	Ip string

	Port string

	Anonymous string

	Protocol string

	Location string
}

func (proxy *Proxy) String() string {
	return fmt.Sprintf("Proxy{ip=%s, port=%s}", proxy.Ip, proxy.Port)
}
