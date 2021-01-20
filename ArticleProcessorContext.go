package IrisAPIs

import (
	"bytes"
	"github.com/pkg/errors"
)

type ArticleProcessorService interface {
	Transform(article string, param ProcessParameters) (string, error)
}

type ProcessParameters struct {
	BytesPerLine int
}

type ArticleProcessorContext struct {
	Param ProcessParameters
}

func NewArticleProcessorContext(param ProcessParameters) (*ArticleProcessorContext, error) {
	if param.BytesPerLine < 5 {
		return nil, errors.New("BytePerLine must greater then 5")
	}
	return &ArticleProcessorContext{Param: param}, nil
}

func (a *ArticleProcessorContext) Transform(article string) (string, error) {
	var buffer bytes.Buffer
	index := 0
	for _, v := range article {

		//fmt.Printf("%d, %c(%U), - %v, index = %d \n", i, v, v, v, index)

		if index >= a.Param.BytesPerLine {
			buffer.WriteString("\n")
			index = 0
		}

		if v == '\n' {
			buffer.WriteString("\n")
			index = 0
			continue
		}

		buffer.WriteString(string(v))

		if v > 256 {
			index += 2
		} else {
			index += 1
		}

	}
	return buffer.String(), nil
}
