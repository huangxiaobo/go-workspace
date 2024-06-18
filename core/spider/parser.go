package spider


type Parser interface {
	Name()string
	Parse(content string) error
}

