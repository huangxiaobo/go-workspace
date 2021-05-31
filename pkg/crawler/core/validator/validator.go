package validator

import (
	"log"
	utils2 "money/pkg/crawler/core/utils"
)
import "github.com/kirinlabs/HttpRequest"

// 检查 proxy 是否是有效

type Validator struct {
}

func validate(proxy utils2.Proxy) bool {
	response, err := HttpRequest.Get("https://www.baidu.com")
	if err != nil {
		log.Fatal(err)
		return false
	}

	log.Printf("response: %v", response)
	return true
}
