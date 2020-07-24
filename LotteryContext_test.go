package IrisAPIs

import (
	"fmt"
	"github.com/jarcoal/httpmock"
	"io/ioutil"
	"testing"
)



func TestLotteryContext_Fetch(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	content, err := ioutil.ReadFile("./test_resources/lottery.html")
	if err != nil {
		fmt.Println(err.Error())
	}
	httpmock.RegisterResponder("GET", "https://www.taiwanlottery.com.tw/index_new.aspx", httpmock.NewStringResponder(200, string(content)))

	l := &LotteryContext{}
	result, err := l.Fetch()
	fmt.Printf("SuperLottery : %+v\n", result.SuperLotto638Result)
}