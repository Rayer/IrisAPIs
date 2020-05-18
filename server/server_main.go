package main

import (
	"IrisAPIs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/moogar0880/problems"
	"net/http"
)

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	chatbot := IrisAPIs.NewChatbotContext()

	r.NoRoute(func(c *gin.Context) {
		err404 := problems.NewStatusProblem(404)
		err404.Detail = "No such route!"
		c.JSON(404, err404)
	})

	r.NoMethod(func(c *gin.Context) {
		err404 := problems.NewStatusProblem(404)
		err404.Detail = "No such method!"
		c.JSON(404, err404)
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello world!",
		})
	})

	r.GET("/currency", func(c *gin.Context) {
		result, err := IrisAPIs.GetMostRecentCurrencyDataRaw()
		if err != nil {
			err500 := problems.NewDetailedProblem(500, err.Error())
			c.JSON(500, err500)
			return
		}

		c.Data(http.StatusOK, "application/json", []byte(result))
	})

	r.GET("/ip2nation", func(c *gin.Context) {
		ipAddr := c.Query("ip")
		if ipAddr == "" {
			err400 := problems.NewDetailedProblem(400, "No query parameter : ip")
			c.JSON(400, err400)
			return
		}
		res, err := IrisAPIs.GetIPNation(ipAddr)
		if err != nil {
			err500 := problems.NewDetailedProblem(500, err.Error())
			c.JSON(500, err500)
			return
		}
		c.JSON(200, res)
	})

	r.POST("/chatbot", func(c *gin.Context) {
		var conv IrisAPIs.ChatbotConversion
		err := c.BindJSON(&conv)

		utx, _ := chatbot.GetUserContext(conv.User)

		prompt, keywords_v, keywords_iv, err := utx.RenderMessageWithDetail()
		str, err := utx.HandleMessage(conv.Input)
		next, err := utx.RenderMessage()
		c.JSON(http.StatusOK, gin.H{
			"prompt":           prompt,
			"keywords":         keywords_v,
			"invalid_keywords": keywords_iv,
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

	IrisAPIs.CurrencySyncRoutine()
	r.Run()
}
