package main

import (
	"github.com/gin-gonic/gin"
	"github.com/moogar0880/problems"
	"net/http"
)

func (c *Controller) GetCurrencyRaw(ctx *gin.Context) {
	result, err := c.CurrencyContext.GetMostRecentCurrencyDataRaw()
	if err != nil {
		err500 := problems.NewDetailedProblem(http.StatusInternalServerError, err.Error())
		ctx.JSON(500, err500)
		return
	}

	ctx.Data(http.StatusOK, "application/json", []byte(result))
}

func (c *Controller) ConvertCurrency(ctx *gin.Context) {
	type payload struct {
		From   string  `json:"from"`
		To     string  `json:"to"`
		Amount float64 `json:"amount"`
		Result float64 `json:"result"`
	}

	var conv payload
	err := ctx.BindJSON(&conv)
	if err != nil {
		err500 := problems.NewDetailedProblem(http.StatusInternalServerError, err.Error())
		ctx.JSON(500, err500)
		return
	}

	if conv.From == "" || conv.To == "" {
		err400 := problems.NewDetailedProblem(http.StatusBadRequest, "either from or to is null!")
		ctx.JSON(400, err400)
		return
	}

	result, err := c.CurrencyContext.Convert(conv.From, conv.To, conv.Amount)
	if err != nil {
		err500 := problems.NewDetailedProblem(http.StatusInternalServerError, err.Error())
		ctx.JSON(500, err500)
		return
	}

	conv.Result = result
	ctx.JSON(200, conv)
}
