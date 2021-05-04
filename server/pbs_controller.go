package main

import (
	"github.com/gin-gonic/gin"
	"github.com/moogar0880/problems"
	"net/http"
	"strconv"
	"time"
)

type GetRecentPBSDataResponse struct {
	ID     string   `json:"id,omitempty"`
	Events []string `json:"events,omitempty"`
}

func (c *Controller) GetRecentPBSData(ctx *gin.Context) {
	p, _ := ctx.GetQuery("period")
	period, err := strconv.Atoi(p)
	if err != nil {
		period = 3600
	}
	result, err := c.PbsTrafficDataService.GetHistory(ctx, time.Duration(period)*time.Second)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, problems.NewDetailedProblem(http.StatusInternalServerError, err.Error()))
		return
	}
	ret := make([]GetRecentPBSDataResponse, 0)
	for k, v := range result {

	}

}
