package main

import (
	"IrisAPIs"
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
}

type GenericResultResponse struct {
	Result bool `json:"result"`
}

func NewController(config *IrisAPIs.Configuration) (*Controller, error) {
	db, err := IrisAPIs.NewDatabaseContext(config.ConnectionString, true, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Error initializing controller!")
	}
	return &Controller{
		ChatBotService:  IrisAPIs.NewChatbotContext(),
		CurrencyService: IrisAPIs.NewCurrencyContextWithConfig(config, db),
		DatabaseContext: db,
		IpNationService: IrisAPIs.NewIpNationContext(db),
		ApiKeyService:   IrisAPIs.NewApiKeyService(db),
		ServiceMgmt: func() IrisAPIs.ServiceManagement {
			service := IrisAPIs.NewServiceManagement()
			_ = service.RegisterPresetServices()
			return service
		}(),
		ArticleProcessorService: IrisAPIs.NewArticleProcessorContext(),
		BuildInfoService:        IrisAPIs.NewBuildInfoService(),
		PbsTrafficDataService:   IrisAPIs.NewPbsTrafficDataService(db),
	}, nil
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
