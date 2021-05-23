package main

import (
	"github.com/gin-gonic/gin"
	"github.com/moogar0880/problems"
	"net/http"
	"strconv"
	"time"
)

type SinglePBSEventInfo struct {
	Time    time.Time `json:"time"`
	Message string    `json:"message,omitempty"`
}

type GetRecentPBSDataResponse struct {
	ID     string               `json:"id,omitempty"`
	Events []SinglePBSEventInfo `json:"events,omitempty"`
}

// GetRecentPBSData godoc
// @Summary Get most recent PBS data
// @Description Get most PBS data, default value is 3600 seconds
// @Tags PBS
// @Accept json
// @Produce json
// @Param period query string false "Period, in seconds(default = 3600)"
// @Success 200 {array} GetRecentPBSDataResponse
// @Failure 400 {object} problems.DefaultProblem
// @Router /pbs/recent [get]
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
		//sort.Slice(v, func(i, j int) bool {
		//	return v[i].LastUpdateTimestamp.Unix() < v[j].LastUpdateTimestamp.Unix()
		//})
		events := make([]SinglePBSEventInfo, 0)
		for _, i := range v {
			events = append(events, SinglePBSEventInfo{
				Time:    *i.LastUpdateTimestamp,
				Message: *i.Information,
			})
		}
		ret = append(ret, GetRecentPBSDataResponse{
			ID:     k,
			Events: events,
		})
	}
	ctx.JSON(http.StatusOK, ret)
}
