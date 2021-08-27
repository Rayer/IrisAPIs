package main

import (
	"IrisAPIs"
	IrisAPIsGRPC "IrisAPIs/grpc"
	"IrisAPIs/server/docs"
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"strings"
	"sync"
	"time"
)

// @title Iris Node Mainframe API
// @version 1.0
// @description This is support APIs for Iris Node
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @basePath /
// @schemes http https

// @securityDefinitions.apikey ApiKeyAuth
// @in query
// @name apiKey

// @in query
// @name apiKey
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host api.rayer.idv.tw
//go:generate go get -u github.com/golang/mock/mockgen@v1.4.4
//go:generate go get -u github.com/swaggo/swag/cmd/swag@v1.6.7
//go:generate ${GOPATH}/bin/swag init -g server_main.go
func main() {

	r := gin.Default()
	r.Use(cors.Default())

	config := &IrisAPIs.Configuration{}
	var log *logrus.Logger
	var controller *Controller
	config.OnFinishedLoadConfig = func(config *IrisAPIs.Configuration) {
		//Init logger
		log, controller = initWithConfig(config, log, r, controller)
	}

	err := config.LoadConfiguration()
	if err != nil {
		panic(err.Error())
	}

	apiKeyManager := NewApiKeyValidator(controller.ApiKeyService, config.EnforceApiKey)
	r.Use(InjectLoggerMiddleware(log))
	r.Use(apiKeyManager.GetMiddleware())

	_ = setupRouter(NewAKWrappedEngine(r, apiKeyManager), controller)

	//Check other services
	ret := controller.ServiceMgmt.CheckAllServerStatus(context.Background())
	format := "%38s %20s %16s %8s %40s\n"
	fmt.Printf(format, "ID", "Name", "Type", "Status", "Message")
	for _, status := range ret {
		fmt.Printf(format, status.ID, status.Name, status.ServiceType, status.Status, status.Message)
	}

	// GRPC server can't reload while configuration changed
	grpc := new(IrisAPIsGRPC.GRPCServerRoutine)
	grpc.RunDetach(context.Background(), config)

	err = r.Run()
	if err != nil {
		panic(err.Error())
	}
}

var once sync.Once
var cancel context.CancelFunc

func initWithConfig(config *IrisAPIs.Configuration, log *logrus.Logger, r *gin.Engine, controller *Controller) (*logrus.Logger, *Controller) {
	log = SetupLogger(config)
	log.Debugf("Configuration : %+v", config)

	if cancel != nil {
		cancel()
	}

	once.Do(func() {
		//Swagger initialization
		_, host, err := config.SplitSchemeAndHost()
		if err != nil {
			panic(err)
		}
		swaggerUrl := ginSwagger.URL(config.Host + "/swagger/doc.json")
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerUrl))

		docs.SwaggerInfo.Host = host
		//Only localhost supports http, for others, force to https
		if strings.HasPrefix(host, "localhost") || strings.HasPrefix(host, "127") {
			docs.SwaggerInfo.Schemes = []string{"http"}
		} else {
			docs.SwaggerInfo.Schemes = []string{"https"}
		}
		controller, err = NewController(config)
		if err != nil {
			panic(err.Error())
		}
	})

	//Run daemon threads
	//TODO: Implement cancel
	var ctx context.Context
	ctx, cancel = context.WithCancel(context.Background())
	if config.CurrencyUpdateRoutine > 0 {
		//TODO: Implement update timer
		controller.CurrencyService.CurrencySyncRoutine(ctx)
		log.Infof("Starting Currency Sync Routine, update for every %d seconds", config.CurrencyUpdateRoutine)
	} else {
		log.Infof("Currency Update Routine is disabled.")
	}
	if config.PBSUpdateRoutine > 0 {
		controller.PbsTrafficDataService.ScheduledRoutine(ctx, time.Duration(config.PBSUpdateRoutine)*time.Second)
		log.Infof("Starting PBS Sync Routine, update for every %d seconds", config.PBSUpdateRoutine)
	} else {
		log.Infof("PBS Update Routine is disabled.")
	}
	return log, controller
}

func setupRouter(wrapped *AKWrappedEngine, controller *Controller) error {

	wrapped.NoRoute(controller.NoRouteHandler)
	wrapped.NoMethod(controller.NoMethodHandler)
	wrapped.GET("/ping", controller.PingHandler)

	system := wrapped.Group("/service")
	{
		system.GET("", IrisAPIs.ApiKeyNotPresented, controller.GetServiceStatus)
		system.GET("/:id", IrisAPIs.ApiKeyNotPresented, controller.GetServiceStatusById)
		system.GET("/:id/logs", IrisAPIs.ApiKeyNotPresented, controller.GetServiceLogs)
	}

	apiKey := wrapped.Group("/apiKey")
	{
		apiKey.POST("", IrisAPIs.ApiKeyPrivileged, controller.IssueApiKey)
		apiKey.GET("", IrisAPIs.ApiKeyPrivileged, controller.GetAllKeys)
		apiKey.GET("/:id/usage", IrisAPIs.ApiKeyPrivileged, controller.GetApiUsage)
		apiKey.GET("/:id", IrisAPIs.ApiKeyPrivileged, controller.GetKey)
	}

	currency := wrapped.Group("/currency")
	{
		currency.GET("", IrisAPIs.ApiKeyNormal, controller.GetCurrencyRaw)
		currency.GET("/sync", IrisAPIs.ApiKeyPrivileged, controller.SyncData)
		currency.POST("", IrisAPIs.ApiKeyNotPresented, controller.ConvertCurrency)
	}

	ipNation := wrapped.Group("/ip2nation")
	{
		ipNation.GET("", IrisAPIs.ApiKeyNormal, controller.IpToNation)
		ipNation.POST("/bulk", IrisAPIs.ApiKeyNormal, controller.IpToNationBulk)
		ipNation.GET("/myip", IrisAPIs.ApiKeyNotPresented, controller.IpToNationMyIP)
	}

	chatbot := wrapped.Group("/chatbot")
	{
		chatbot.POST("", IrisAPIs.ApiKeyNormal, controller.ChatBotReact)
		chatbot.DELETE("/:user", IrisAPIs.ApiKeyPrivileged, controller.ChatBotResetUser)
	}

	articleProcess := wrapped.Group("/article_process")
	{
		articleProcess.POST("", IrisAPIs.ApiKeyNotPresented, controller.TransformArticle)
	}

	pbs := wrapped.Group("/pbs")
	{
		pbs.GET("/recent", IrisAPIs.ApiKeyNotPresented, controller.GetRecentPBSData)
	}

	gLogger.Info("Listing privilege endpoints : ")
	privilegeEndpoints := wrapped.GetPrivilegeMap()
	for path, level := range privilegeEndpoints {
		gLogger.Infof("%s(%#v)", path, level)
	}

	return nil
}
