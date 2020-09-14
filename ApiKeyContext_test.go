package IrisAPIs

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ApiKeyContextTestSuite struct {
	suite.Suite
	db      *DatabaseContext
	context *ApiKeyContext
}

func (c *ApiKeyContextTestSuite) SetupTest() {
	c.db, _ = NewDatabaseContext("acc:12qw34er@tcp(node.rayer.idv.tw:3306)/apps_test?charset=utf8&loc=Asia%2FTaipei&parseTime=true", true)
	c.context = &ApiKeyContext{DB: c.db.DbObject}
}

func TestApiKeyContextTestSuite(t *testing.T) {
	suite.Run(t, new(ApiKeyContextTestSuite))
}

func (c *ApiKeyContextTestSuite) TestApiKeyContext_IssueApiKey() {
	fmt.Println(c.context.IssueApiKey("TestApplication", true, true))

}

func (c *ApiKeyContextTestSuite) TestValidateApiKey() {
	fmt.Println(c.context.ValidateApiKey("12345anvc", HEADER))
	fmt.Println(c.context.ValidateApiKey("WxlP7RgaUqE7Q8so", QUERY_STRING))
}
