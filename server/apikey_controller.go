package main

import (
	"github.com/gin-gonic/gin"
	"github.com/moogar0880/problems"
	"net/http"
)

type IssueApiKeyPayload struct {
	Application     string `json:"application"`
	UseInHeader     bool   `json:"use_in_header"`
	UseInQueryParam bool   `json:"use_in_query_param"`
}

type IssueApiKeyResponse struct {
	Key string `json:"key"`
}

// IssueApiKey godoc
// @Summary Issue an API Key
// @Description Issue an API Key to user, this endpoint requires privileges
// @Tags ApiKey
// @Accept json
// @Produce json
// @param input body IssueApiKeyPayload true "Input info"
// @Success 200 {object} IssueApiKeyResponse
// @Failure 400 {object} problems.DefaultProblem
// @Router /apiKey [post]
func (c *Controller) IssueApiKey(ctx *gin.Context) {
	input := &IssueApiKeyPayload{}
	err := ctx.BindJSON(input)
	if err != nil {
		err400 := problems.NewDetailedProblem(http.StatusBadRequest, err.Error())
		ctx.JSON(400, err400)
		return
	}
	key, err := c.ApiKeyService.IssueApiKey(input.Application, input.UseInHeader, input.UseInQueryParam)
	if err != nil {
		err500 := problems.NewDetailedProblem(http.StatusInternalServerError, err.Error())
		ctx.JSON(500, err500)
		return
	}
	ctx.JSON(http.StatusOK, IssueApiKeyResponse{
		Key: key,
	})
}
