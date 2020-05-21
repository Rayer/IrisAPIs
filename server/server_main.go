package main

import (
	"IrisAPIs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/moogar0880/problems"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {

	//Init logger
	log.SetLevel(log.DebugLevel)

	r := gin.Default()
	r.Use(cors.Default())

	config := &IrisAPIs.Configuration{}
	err := config.LoadConfiguration()
	if err != nil {
		panic(err.Error())
	}
	log.Debugf("Configuration : %+v", config)

	db, err := IrisAPIs.NewDatabaseContext(config.ConnectionString, true)
	if err != nil {
		panic(err.Error())
	}
	chatbot := IrisAPIs.NewChatbotContext()

	r.NoRoute(func(c *gin.Context) {
		err404 := problems.NewStatusProblem(http.StatusNotFound)
		err404.Detail = "No such route!"
		c.JSON(404, err404)
	})

	r.NoMethod(func(c *gin.Context) {
		err404 := problems.NewStatusProblem(http.StatusNotFound)
		err404.Detail = "No such method!"
		c.JSON(404, err404)
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello world!",
		})
	})

	//Currency Rate
	currencyContext := IrisAPIs.NewCurrencyContextWithConfig(config, db)

	r.GET("/currency", func(c *gin.Context) {
		result, err := currencyContext.GetMostRecentCurrencyDataRaw()
		if err != nil {
			err500 := problems.NewDetailedProblem(http.StatusInternalServerError, err.Error())
			c.JSON(500, err500)
			return
		}

		c.Data(http.StatusOK, "application/json", []byte(result))
	})

	ipNation := IrisAPIs.NewIpNationContext(db)

	r.GET("/ip2nation", func(c *gin.Context) {
		ipAddr := c.Query("ip")
		if ipAddr == "" {
			err400 := problems.NewDetailedProblem(http.StatusBadRequest, "No query parameter : ip")
			c.JSON(400, err400)
			return
		}
		res, err := ipNation.GetIPNation(ipAddr)
		if err != nil {
			err500 := problems.NewDetailedProblem(http.StatusInternalServerError, err.Error())
			c.JSON(500, err500)
			return
		}
		c.JSON(200, res)
	})

	r.POST("/chatbot", func(c *gin.Context) {
		var conv IrisAPIs.ChatbotConversion
		err := c.BindJSON(&conv)

		utx, _ := chatbot.GetUserContext(conv.User)

		prompt, keywordsV, keywordsIv, err := utx.RenderMessageWithDetail()
		str, err := utx.HandleMessage(conv.Input)
		next, err := utx.RenderMessage()
		c.JSON(http.StatusOK, gin.H{
			"prompt":           prompt,
			"keywords":         keywordsV,
			"invalid_keywords": keywordsIv,
			"message":          str,
			"error":            err,
			"next":             next,
		})
	})

	r.DELETE("/chatbot/:user", func(c *gin.Context) {
		user := c.Param("user")
		chatbot.ExpireUser(user, func() {
			c.JSON(201, gin.H{
				"message": "ok",
			})
		}, func() {
			c.JSON(http.StatusBadRequest, problems.NewDetailedProblem(http.StatusBadRequest, "User "+user+" not found!"))
		})

	})

	//Run daemon threads
	currencyContext.CurrencySyncRoutine()

	err = r.Run()
	if err != nil {
		panic(err.Error())
	}
}
