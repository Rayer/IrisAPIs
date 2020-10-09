package IrisAPIs

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
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
	key, err := c.context.IssueApiKey("TestApplication", true, true)
	if err != nil {
		c.Fail("error while trying issuing apikey", err)
	}
	//This key should able to be validated
	result := c.context.ValidateApiKey(key, ApiKeyLocation(0))
	assert.True(c.T(), result != ApiKeyNotValid)

	//Generate random one and it should not be validated
	assert.True(c.T(), c.context.ValidateApiKey("abcd1234", ApiKeyLocation(0)) == ApiKeyNotValid)
}
