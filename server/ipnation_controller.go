package main

import (
	"github.com/gin-gonic/gin"
	"github.com/moogar0880/problems"
	"net/http"
)

func (c *Controller) IpToNation(ctx *gin.Context) {
	ipAddr := ctx.Query("ip")
	if ipAddr == "" {
		err400 := problems.NewDetailedProblem(http.StatusBadRequest, "No query parameter : ip")
		ctx.JSON(400, err400)
		return
	}
	res, err := c.IpNationContext.GetIPNation(ipAddr)
	if err != nil {
		err500 := problems.NewDetailedProblem(http.StatusInternalServerError, err.Error())
		ctx.JSON(500, err500)
		return
	}
	ctx.JSON(200, res)
}
