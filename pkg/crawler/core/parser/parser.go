package parser

import (
	utils2 "money/pkg/crawler/core/utils"
)

type Parser interface {
	Parse(content string) []utils2.Proxy
}

func Factory(domain utils2.DomainType) Parser {

	switch domain {
	case utils2.KUAIDAILI:
		return &Kuaidaili{}
	case utils2.XICI:
		return &Xicidaili{}
	}

	return nil

}
