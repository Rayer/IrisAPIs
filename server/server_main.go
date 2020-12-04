package main

import (
	"IrisAPIs"
	"IrisAPIsServer/docs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
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
	log.Debugf("Configuration : %+v", config)

	//Init logger
	log.SetLevel(log.Level(config.LogLevel))
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

	r.NoRoute(controller.NoRouteHandler)
	r.NoMethod(controller.NoMethodHandler)

	r.GET("/ping", controller.PingHandler)

	wrapped := NewAKWrappedEngine(r, apiKeyManager)

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
	}

	chatbot := wrapped.Group("/chatbot")
	{
		chatbot.POST("", IrisAPIs.ApiKeyNormal, controller.ChatBotReact)
		chatbot.DELETE("/:user", IrisAPIs.ApiKeyPrivileged, controller.ChatBotResetUser)
	}

	//Run daemon threads
	controller.CurrencyService.CurrencySyncRoutine()

	log.Info("Listing privilege endpoints : ")
	privilegeEndpoints := wrapped.GetPrivilegeMap()
	for path, level := range privilegeEndpoints {
		log.Infof("%s(%#v)", path, level)
	}

	err = r.Run()
	if err != nil {
		panic(err.Error())
	}
}
