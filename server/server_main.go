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
	err := config.LoadConfiguration()
	if err != nil {
		panic(err.Error())
	}
	log := gLogger
	log.SetLevel(logrus.Level(config.LogLevel))

	log.Debugf("Configuration : %+v", config)

	//Init logger
	//Swagger initialization
	_, host, err := config.SplitSchemeAndHost()
	if err != nil {
		panic(err)
	}
	swaggerUrl := ginSwagger.URL(config.Host + "/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerUrl))

	docs.SwaggerInfo.Host = host

	controller, err := NewController(config)
	if err != nil {
		panic(err.Error())
	}
	apiKeyManager := NewApiKeyValidator(controller.ApiKeyService, config.EnforceApiKey)
	r.Use(apiKeyManager.GetMiddleware())
	r.Use(LoggerMiddleware(gLogger))

	_ = setupRouter(NewAKWrappedEngine(r, apiKeyManager), controller)

	//Run daemon threads
	controller.CurrencyService.CurrencySyncRoutine()

	//Check other services
	ret := controller.ServiceMgmt.CheckAllServerStatus()
	format := "%38s %20s %16s %8s %40s\n"
	fmt.Printf(format, "ID", "Name", "Type", "Status", "Message")
	for _, status := range ret {
		fmt.Printf(format, status.ID, status.Name, status.ServiceType, status.Status, status.Message)
	}

	grpc := new(IrisAPIsGRPC.GRPCServerRoutine)
	grpc.RunDetach(context.TODO(), config)

	err = r.Run()
	if err != nil {
		panic(err.Error())
	}
}

func setupRouter(wrapped *AKWrappedEngine, controller *Controller) error {

	wrapped.NoRoute(controller.NoRouteHandler)
	wrapped.NoMethod(controller.NoMethodHandler)
	wrapped.GET("/ping", controller.PingHandler)

	system := wrapped.Group("/service")
	{
		system.GET("", IrisAPIs.ApiKeyNotPresented, controller.GetServiceStatus)
		system.GET("/:id", IrisAPIs.ApiKeyNotPresented, controller.GetServiceStatusById)
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

	gLogger.Info("Listing privilege endpoints : ")
	privilegeEndpoints := wrapped.GetPrivilegeMap()
	for path, level := range privilegeEndpoints {
		gLogger.Infof("%s(%#v)", path, level)
	}

	return nil
}
