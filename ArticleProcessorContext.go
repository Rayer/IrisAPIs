package IrisAPIs

import (
	"bytes"
	"github.com/pkg/errors"
)

type ArticleProcessorService interface {
	Transform(param ProcessParameters, article string) (string, error)
	validateProcessParameters(param ProcessParameters) error
}

type ProcessParameters struct {
	BytesPerLine int
}

type ArticleProcessorContext struct {
}

func (a *ArticleProcessorContext) validateProcessParameters(param ProcessParameters) error {
	if param.BytesPerLine < 5 {
		return errors.New("BytesPerLine requires at least 5")
	}
	return nil
}

func NewArticleProcessorContext() *ArticleProcessorContext {
	return &ArticleProcessorContext{}
}

func (a *ArticleProcessorContext) Transform(param ProcessParameters, article string) (string, error) {
	err := a.validateProcessParameters(param)
	if err != nil {
		return "", err
	}
	var buffer bytes.Buffer
	index := 0
	for _, v := range article {

		//fmt.Printf("%d, %c(%U), - %v, index = %d \n", i, v, v, v, index)

		if index >= param.BytesPerLine {
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
