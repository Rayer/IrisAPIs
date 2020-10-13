package main

import (
	"github.com/gin-gonic/gin"
	"github.com/moogar0880/problems"
	"net/http"
)

type CurrencyConvert struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
	Result float64 `json:"result"`
}

// GetCurrencyRaw godoc
// @Summary Get most recent raw data
// @Description Get most recent raw data fetching from fixer.io
// @Tags Currency
// @Accept json
// @Produce json
// @Param apiKey query string true "API Key"
// @Success 200 {string} string "...Data from source"
// @Failure 400 {object} problems.DefaultProblem
// @Router /currency [get]
func (c *Controller) GetCurrencyRaw(ctx *gin.Context) {
	result, err := c.CurrencyService.GetMostRecentCurrencyDataRaw()
	if err != nil {
		err500 := problems.NewDetailedProblem(http.StatusInternalServerError, err.Error())
		ctx.JSON(500, err500)
		return
	}

	ctx.Data(http.StatusOK, "application/json", []byte(result))
}

// ConvertCurrency godoc
// @Summary Convert currency
// @Description Convert currency from most recent data
// @Tags Currency
// @param input body CurrencyConvert true "Currency Convert info"
// @Accept json
// @Produce json
// @Success 200 {object} CurrencyConvert
// @Failure 400 {object} problems.DefaultProblem
// @Router /currency/convert [post]
func (c *Controller) ConvertCurrency(ctx *gin.Context) {

	var conv CurrencyConvert
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

	result, err := c.CurrencyService.Convert(conv.From, conv.To, conv.Amount)
	if err != nil {
		err500 := problems.NewDetailedProblem(http.StatusInternalServerError, err.Error())
		ctx.JSON(500, err500)
		return
	}

	conv.Result = result
	ctx.JSON(200, conv)
}
