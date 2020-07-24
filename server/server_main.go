package main

import (
	"IrisAPIs"
	"IrisAPIsServer/docs"
	_ "IrisAPIsServer/docs"
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

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host api.rayer.idv.tw
//go:generate go get -u github.com/swaggo/swag/cmd/swag
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

	r.NoRoute(controller.NoRouteHandler)
	r.NoMethod(controller.NoMethodHandler)
	r.GET("/ping", controller.PingHandler)

	//system := r.Group("/system")
	//{
	//
	//}

	currency := r.Group("/currency")
	{
		currency.GET("", controller.GetCurrencyRaw)
		currency.POST("/convert", controller.ConvertCurrency)
	}

	ipNation := r.Group("/ip2nation")
	{
		ipNation.GET("", controller.IpToNation)
		ipNation.POST("/bulk", controller.IpToNationBulk)
	}

	chatbot := r.Group("/chatbot")
	{
		chatbot.POST("", controller.ChatBotReact)
		chatbot.DELETE("/:user", controller.ChatBotResetUser)
	}

	//Run daemon threads
	//IrisAPIs.NewCurrencyContextWithConfig(config, controller.DatabaseContext).CurrencySyncRoutine()
	controller.CurrencyContext.CurrencySyncRoutine()

	err = r.Run()
	if err != nil {
		panic(err.Error())
	}
}
