package IrisAPIs

import (
	"fmt"
	"testing"
)

func TestArticleProcessorContext_Transform(t *testing.T) {
	a := ArticleProcessorContext{}
	s, _ := a.Transform(`要實作介面，你只需要實作介面中宣告的所有方法。
Go 的介面是隱含實作的
與 Java 等其他語言不同，你不需要使用如 implements 關鍵字之類的方法來明確指定一種型別來實作介面。
以下是兩種實作 Shape 介面的 struct 型別：`, ProcessParameters{
		BytesPerLine: 25,
	})
	fmt.Println(s)
}
