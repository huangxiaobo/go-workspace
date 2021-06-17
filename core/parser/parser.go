package parser

type Parser interface {
	Parse(content string) error
}

const ParserZhihu = "zhihu"

func Factory(name string) Parser {

	switch name {
	case ParserZhihu:
		return &ZhihuPeople{}
	}

	return nil

}
