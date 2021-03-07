package parser

import "money/core/utils"

type Parser interface {
	Parse(content string) [] utils.Proxy
}

func Factory(domain utils.DomainType) Parser {

	switch domain {
	case utils.KUAIDAILI:
		return &Kuaidaili{}
	case utils.XICI:
		return &Xicidaili{}
	}

	return nil

}
