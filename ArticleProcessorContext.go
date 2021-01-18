package IrisAPIs

import (
	"bytes"
)

type ArticleProcessorService interface {
	Transform(article string, param ProcessParameters) (string, error)
}

type ProcessParameters struct {
	BytesPerLine int
}

type ArticleProcessorContext struct {
}

func (a *ArticleProcessorContext) Transform(article string, param ProcessParameters) (string, error) {
	var buffer bytes.Buffer
	index := 0
	for _, v := range article {

		//fmt.Printf("%d, %c(%U), - %v, index = %d \n", i, v, v, v, index)
		buffer.WriteString(string(v))

		if index >= param.BytesPerLine {
			buffer.WriteString("\n")
			index = 0
		}

		if v == '\n' {
			buffer.WriteString("\n")
			index = 0
			continue
		}

		if v > 256 {
			index += 2
		} else {
			index += 1
		}

	}
	return buffer.String(), nil
}
