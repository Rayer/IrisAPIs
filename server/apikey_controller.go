package main

import (
	"github.com/gin-gonic/gin"
	"github.com/moogar0880/problems"
	"net/http"
	"strconv"
	"time"
)

type IssueApiKeyPayload struct {
	Application     string `json:"application"`
	UseInHeader     bool   `json:"use_in_header"`
	UseInQueryParam bool   `json:"use_in_query_param"`
}

type IssueApiKeyResponse struct {
	Key string `json:"key"`
}

type ApiKeyDetail struct {
	ApiKeyBrief
	IssueBy     string
	Application string
}

type ApiKeyBrief struct {
	Id         int    `json:"id"`
	Key        string `json:"key"`
	Privileged bool   `json:"privileged"`
}

type AccessRecord struct {
	Path string    `json:"path"`
	Ip   string    `json:"ip"`
	Time time.Time `json:"time"`
}

type ApiKeyUsage struct {
	Id     int
	Access []AccessRecord
}

// IssueApiKey godoc
// @Summary Issue an API Key
// @Description Issue an API Key to user, this endpoint requires privileges
// @Tags ApiKey
// @Accept json
// @Produce json
// @param input body IssueApiKeyPayload true "Input info"
// @Param apiKey query string true "API Key(Privileged)"
// @Security ApiKeyAuth
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
	key, err := c.ApiKeyService.IssueApiKey(input.Application, input.UseInHeader, input.UseInQueryParam, "auto", false)
	if err != nil {
		err500 := problems.NewDetailedProblem(http.StatusInternalServerError, err.Error())
		ctx.JSON(500, err500)
		return
	}
	ctx.JSON(http.StatusOK, IssueApiKeyResponse{
		Key: key,
	})
}

// @Summary Get API Key list
// @Description Get current api keys
// @Tags ApiKey
// @Accept json
// @Produce json
// @Param apiKey query string true "API Key(Privileged)"
// @Security ApiKeyAuth
// @Success 200 {array} ApiKeyBrief
// @Failure 400 {object} problems.DefaultProblem
// @Router /apiKey [get]
func (c *Controller) GetAllKeys(ctx *gin.Context) {
	entities, err := c.ApiKeyService.GetAllKeys()
	if err != nil {
		err500 := problems.NewDetailedProblem(http.StatusInternalServerError, err.Error())
		ctx.JSON(500, err500)
		return
	}

	ret := make([]ApiKeyBrief, 0)

	for _, entity := range entities {
		ret = append(ret, ApiKeyBrief{
			Id:         *entity.Id,
			Key:        *entity.Key,
			Privileged: *entity.Privileged,
		})
	}

	ctx.JSON(http.StatusOK, ret)
}

// @Summary Get API Key detail
// @Description Get destinated API Key detail
// @Tags ApiKey
// @Accept json
// @Produce json
// @Param id path integer true "Api Key ID"
// @Param apiKey query string true "API Key(Privileged)"
// @Security ApiKeyAuth
// @Success 200 {array} ApiKeyDetail
// @Failure 400 {object} problems.DefaultProblem
// @Router /apiKey/{id} [get]
func (c *Controller) GetKey(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		err400 := problems.NewDetailedProblem(http.StatusBadRequest, "Bad ID")
		ctx.JSON(http.StatusBadRequest, err400)
		return
	}

	entity, err := c.ApiKeyService.GetKey(id)
	if err != nil {
		err500 := problems.NewDetailedProblem(http.StatusInternalServerError, err.Error())
		ctx.JSON(http.StatusInternalServerError, err500)
		return
	}

	if entity == nil {
		err404 := problems.NewDetailedProblem(http.StatusNotFound, "ID not found")
		ctx.JSON(http.StatusNotFound, err404)
		return
	}

	ctx.JSON(http.StatusOK, ApiKeyDetail{
		ApiKeyBrief: ApiKeyBrief{
			Id:         *entity.Id,
			Key:        *entity.Key,
			Privileged: *entity.Privileged,
		},
		IssueBy:     *entity.Issuer,
		Application: *entity.Application,
	})
}

// IssueApiKey godoc
// @Summary Get API Usages
// @Description Get API Usages, can pass timestamp into thee
// @Tags ApiKey
// @Accept json
// @Produce json
// @Param id path integer true "Api Key ID"
// @Param from query integer false "From(timestamp)"
// @Param to query integer false "To(timestamp)"
// @Param apiKey query string true "API Key(Privileged)"
// @Security ApiKeyAuth
// @Success 200 {object} ApiKeyUsage
// @Failure 400 {object} problems.DefaultProblem
// @Router /apiKey/{id}/usage [get]
func (c *Controller) GetApiUsage(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		err400 := problems.NewDetailedProblem(http.StatusBadRequest, "Bad ID")
		ctx.JSON(http.StatusBadRequest, err400)
		return
	}

	from, _ := strconv.ParseInt(ctx.Query("from"), 10, 64)
	to, _ := strconv.ParseInt(ctx.Query("to"), 10, 64)

	var fromT *time.Time
	if from > 0 {
		t := time.Unix(from, 0)
		fromT = &t
	}
	var toT *time.Time
	if to > 0 {
		t := time.Unix(to, 0)
		toT = &t
	}

	records, err := c.ApiKeyService.GetKeyUsage(id, fromT, toT)

	if err != nil {
		err500 := problems.NewDetailedProblem(http.StatusInternalServerError, err.Error())
		ctx.JSON(500, err500)
		return
	}

	ret := make([]AccessRecord, 0)
	for _, v := range records {
		ret = append(ret, AccessRecord{
			Path: *v.Fullpath,
			Ip:   *v.Ip,
			Time: *v.Timestamp,
		})
	}

	ctx.JSON(http.StatusOK, ApiKeyUsage{
		Id:     id,
		Access: ret,
	})

}
