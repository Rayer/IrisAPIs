package main

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ControllerTestSuite struct {
	suite.Suite
	gin              *gin.Context
	controller       *Controller
	responseRecorder *httptest.ResponseRecorder
}

func TestControllerTestSuite(t *testing.T) {
	suite.Run(t, new(ControllerTestSuite))
}

func (c *ControllerTestSuite) SetupTest() {
	c.responseRecorder = httptest.NewRecorder()
	c.gin, _ = gin.CreateTestContext(c.responseRecorder)
	//We doesn't need services in controller, just test controller itself
	c.controller = &Controller{}
}

func (c *ControllerTestSuite) TestController_NoMethodHandler() {
	c.controller.NoMethodHandler(c.gin)
	assert.Equal(c.T(), http.StatusNotFound, c.responseRecorder.Code)
}

func (c *ControllerTestSuite) TestController_NoRouteHandler() {
	c.controller.NoRouteHandler(c.gin)
	assert.Equal(c.T(), http.StatusNotFound, c.responseRecorder.Code)
}

func (c *ControllerTestSuite) TestController_PingHandler() {
	c.controller.PingHandler(c.gin)
	assert.Equal(c.T(), http.StatusOK, c.responseRecorder.Code)
}
