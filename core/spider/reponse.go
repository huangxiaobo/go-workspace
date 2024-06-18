package spider

import "io"

type Reponse struct {
	Body io.ReadCloser

	CSS(string)
}


