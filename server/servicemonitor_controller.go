package main

import (
	"github.com/docker/distribution/uuid"
	"github.com/gin-gonic/gin"
	"github.com/moogar0880/problems"
	"net/http"
)

type GetServiceStatusByIdResponse struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// GetServiceStatus godoc
// @Summary Get all service status
// @Description Get all service status upon registered ones
// @Tags System
// @Produce json
// @Success 200 {array} GetServiceStatusByIdResponse
// @Router /service [get]
func (c *Controller) GetServiceStatus(ctx *gin.Context) {
	ret := make([]GetServiceStatusByIdResponse, 0)
	for _, stat := range c.ServiceMgmt.CheckAllServerStatus(ctx) {
		ret = append(ret, GetServiceStatusByIdResponse{
			Id:      stat.ID.String(),
			Name:    stat.Name,
			Type:    stat.ServiceType,
			Status:  string(stat.Status),
			Message: stat.Message,
		})
	}
	ctx.JSON(http.StatusOK, ret)
}

// GetServiceStatusById godoc
// @Summary Get service status
// @Description Get service status with specified ID
// @Tags System
// @Param id path string true "Service ID"
// @Produce json
// @Success 200 {object} GetServiceStatusByIdResponse
// @Failure 400 {object} problems.DefaultProblem
// @Router /service/{id} [get]
func (c *Controller) GetServiceStatusById(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		err400 := problems.NewDetailedProblem(http.StatusBadRequest, "error parsing service id")
		ctx.JSON(400, err400)
		return
	}

	stat, err := c.ServiceMgmt.CheckServerStatus(ctx, id)
	if err != nil {
		err404 := problems.NewDetailedProblem(http.StatusNotFound, "no such service bound with this id")
		ctx.JSON(http.StatusNotFound, err404)
		return
	}

	ctx.JSON(http.StatusOK, GetServiceStatusByIdResponse{
		Id:      stat.ID.String(),
		Name:    stat.Name,
		Type:    stat.ServiceType,
		Status:  string(stat.Status),
		Message: stat.Message,
	})
}

// GetServiceLogs godoc
// @Summary Get service logs
// @Description Get service logs with specified ID
// @Tags System
// @Param id path string true "Service ID"
// @Produce plain
// @Success 200 {string} string "logs here"
// @Failure 400 {object} problems.DefaultProblem
// @Router /service/{id}/logs [get]
func (c *Controller) GetServiceLogs(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		err400 := problems.NewDetailedProblem(http.StatusBadRequest, "error parsing service id")
		ctx.JSON(400, err400)
		return
	}

	logs, err := c.ServiceMgmt.GetLogs(ctx, id)

	if err != nil {
		err400 := problems.NewDetailedProblem(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, err400)
		return
	}

	ctx.String(http.StatusOK, logs)
}
