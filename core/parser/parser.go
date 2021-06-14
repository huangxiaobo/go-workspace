package parser

type Parser interface {
	Parse(content string)
}

func Factory(name string) Parser {

	switch name {
	case "Zhihu":
		return nil
	}

	return nil

}
