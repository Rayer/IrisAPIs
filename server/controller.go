package main

import (
	"IrisAPIs"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/moogar0880/problems"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"time"
)

type SystemDefaultController interface {
	NoRouteHandler(c *gin.Context)
	NoMethodHandler(c *gin.Context)
	PingHandler(c *gin.Context)
}

type Controller struct {
	SystemDefaultController
	ChatBotService          *IrisAPIs.ChatbotContext
	CurrencyService         IrisAPIs.CurrencyService
	DatabaseContext         *IrisAPIs.DatabaseContext
	IpNationService         *IrisAPIs.IpNationContext
	ApiKeyService           IrisAPIs.ApiKeyService
	ServiceMgmt             IrisAPIs.ServiceManagement
	ArticleProcessorService IrisAPIs.ArticleProcessorService
	PbsTrafficDataService   IrisAPIs.PbsTrafficDataService
	BuildInfoService        IrisAPIs.BuildInfoService
	teardownQueue           []IrisAPIs.TeardownableServices
}

type GenericResultResponse struct {
	Result bool `json:"result"`
}

func NewController(config *IrisAPIs.Configuration) (*Controller, error) {
	ret := &Controller{}
	err := ret.ReInitServices(context.TODO(), config)
	return ret, err
}

func (c *Controller) ReInitServices(ctx context.Context, config *IrisAPIs.Configuration) error {
	logger := IrisAPIs.GetLogger(ctx)
	db, err := IrisAPIs.NewDatabaseContext(config.ConnectionString, true, nil)
	if err != nil {
		//If failed to initialize, will stop re-init
		return errors.Wrap(err, "Error initializing database!")
	}

	for _, s := range c.teardownQueue {
		err := s.Teardown()
		if err != nil {
			logger.Warning("%v teardown failed!")
		}
	}

	c.CurrencyService = c.registerService(IrisAPIs.NewCurrencyContextWithConfig(config.FixerIoApiKey,
		config.FixerIoLastFetchSuccessfulPeriod, config.FixerIoLastFetchFailedPeriod, db)).(IrisAPIs.CurrencyService)
	c.ChatBotService = c.registerService(IrisAPIs.NewChatbotContext()).(*IrisAPIs.ChatbotContext)
	c.DatabaseContext = c.registerService(db).(*IrisAPIs.DatabaseContext)
	c.IpNationService = c.registerService(IrisAPIs.NewIpNationContext(db)).(*IrisAPIs.IpNationContext)
	c.ApiKeyService = c.registerService(IrisAPIs.NewApiKeyService(db)).(IrisAPIs.ApiKeyService)
	c.ServiceMgmt = c.registerService(func() IrisAPIs.ServiceManagement {
		service := IrisAPIs.NewServiceManagement()
		_ = service.RegisterPresetServices()
		return service
	}()).(IrisAPIs.ServiceManagement)
	c.ArticleProcessorService = c.registerService(IrisAPIs.NewArticleProcessorContext()).(IrisAPIs.ArticleProcessorService)
	c.BuildInfoService = c.registerService(IrisAPIs.NewBuildInfoService()).(IrisAPIs.BuildInfoService)
	c.PbsTrafficDataService = c.registerService(IrisAPIs.NewPbsTrafficDataService(db)).(IrisAPIs.PbsTrafficDataService)
	return nil
}

func (c *Controller) registerService(service interface{}) interface{} {
	if c.teardownQueue == nil {
		c.teardownQueue = make([]IrisAPIs.TeardownableServices, 0)
	}
	teardownableService, isTeardownable := service.(IrisAPIs.TeardownableServices)

	if isTeardownable {
		c.teardownQueue = append(c.teardownQueue, teardownableService)
	}

	return service
}

func (c *Controller) NoRouteHandler(ctx *gin.Context) {
	err404 := problems.NewStatusProblem(http.StatusNotFound)
	err404.Detail = "No such route!"
	ctx.JSON(404, err404)
}

func (c *Controller) NoMethodHandler(ctx *gin.Context) {
	err404 := problems.NewStatusProblem(http.StatusNotFound)
	err404.Detail = "No such method!"
	ctx.JSON(404, err404)
}

type PingResponse struct {
	Message        string
	Hostname       string
	Timezone       string
	Time           string
	ImageTag       string
	JenkinsUrl     string
	BuildTimestamp int64
}

// PingHandler godoc
// @Summary Ping for Iris health check
// @Tags System
// @Accept json
// @Produce json
// @Success 200 {object} PingResponse
// @Failure 400 {object} problems.DefaultProblem
// @Failure 500 {object} problems.DefaultProblem
// @Router /ping [get]
func (c *Controller) PingHandler(ctx *gin.Context) {
	ctx.Copy()
	hostname, _ := os.Hostname()
	buildInfo := c.BuildInfoService.GetBuildInfo(ctx)
	//_ = os.UserCacheDir()
	ctx.JSON(200, PingResponse{
		Message:        "System alive!!!",
		Hostname:       hostname,
		Timezone:       fmt.Sprint(time.Now().Zone()),
		Time:           fmt.Sprint(time.Now().Format("2006-01-02T15:04:05.000 MST")),
		ImageTag:       buildInfo.ImageTag,
		JenkinsUrl:     buildInfo.JenkinsLink,
		BuildTimestamp: buildInfo.CreateTimestamp,
	})
}
