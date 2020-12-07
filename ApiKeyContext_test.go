package IrisAPIs

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/xormplus/xorm"
	"testing"
	"time"
)

type ApiKeyContextTestSuite struct {
	suite.Suite
	db      *DatabaseContext
	context *ApiKeyContext
}

func (c *ApiKeyContextTestSuite) SetupSuite() {
	c.db, _ = NewTestDatabaseContext()
	c.context = &ApiKeyContext{DB: func() *xorm.Engine {
		if c.db == nil {
			return nil
		}
		return c.db.DbObject
	}()}
}

func (c *ApiKeyContextTestSuite) SetupTest() {
	if c.db == nil {
		c.T().Skip("Skip these tests due to no database available!")
	}
}

func TestApiKeyContextTestSuite(t *testing.T) {
	suite.Run(t, new(ApiKeyContextTestSuite))
}

func (c *ApiKeyContextTestSuite) TestApiKeyContext_IssueApiKey() {
	key, err := c.context.IssueApiKey("TestApplication", true, true, "TestUser", false)
	if err != nil {
		c.Fail("error while trying issuing apikey", err)
	}
	//This key should able to be validated
	result := c.context.ValidateApiKey(key, ApiKeyLocation(0))
	assert.True(c.T(), result != ApiKeyNotValid)

	//Generate random one and it should not be validated
	assert.True(c.T(), c.context.ValidateApiKey("abcd1234", ApiKeyLocation(0)) == ApiKeyNotValid)
}

func (c *ApiKeyContextTestSuite) TestApiKeyContext_GetAllKeys() {
	ret, err := c.context.GetAllKeys()
	if err != nil {
		c.Failf("Error getting keys : %s", err.Error())
		c.Assert()
	}
	for _, r := range ret {
		c.T().Logf("%+v", r)
	}
}

func (c *ApiKeyContextTestSuite) TestApiKeyContext_GetKeyUsage() {
	r1, _ := c.context.GetKeyUsageById(3, nil, nil)
	now := time.Now()
	r2, _ := c.context.GetKeyUsageById(3, nil, &now)

	for _, r := range r1 {
		c.T().Logf("%+v", r)
	}

	for _, r := range r2 {
		c.T().Logf("%+v", r)
	}
}
