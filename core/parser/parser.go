package parser

import "money/core/parser/zhihu"

type Parser interface {
	Parse(content string) error
}

const ParserZhihu = "zhihu"

func Factory(name string) Parser {

	switch name {
	case ParserZhihu:
		return &zhihu.ParserZhihuPeople{}
	}

	return nil
}
