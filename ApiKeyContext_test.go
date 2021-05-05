package IrisAPIs

import (
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type ApiKeyContextTestSuite struct {
	suite.Suite
	db      *DatabaseContext
	context ApiKeyService
}

func (c *ApiKeyContextTestSuite) SetupSuite() {
	c.db, _ = NewTestDatabaseContext()
	c.context = NewApiKeyService(c.db)
	log.SetLevel(log.DebugLevel)
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
	_, result := c.context.ValidateApiKey(key, ApiKeyLocation(0))
	assert.True(c.T(), result != ApiKeyNotValid)

	//Generate random one and it should not be validated
	id, level := c.context.ValidateApiKey("abcd1234", ApiKeyLocation(0))
	assert.Equal(c.T(), ApiKeyNotValid, level)
	assert.Equal(c.T(), -1, id)
}

func (c *ApiKeyContextTestSuite) TestApiKeyContext_GetAllKeys() {
	ret, err := c.context.GetAllKeys()
	if err != nil {
		c.Failf("Error getting keys : %s", err.Error())
		c.Assert()
	}
	//for _, r := range ret {
	//	c.T().Logf("%+v", r)
	//}
	c.T().Logf("Fetched %d keys", len(ret))
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
