package main

import (
	"github.com/Rayer/IrisAPIs"
	"github.com/gin-gonic/gin"
	"github.com/moogar0880/problems"
	"net/http"
)

type TransformArticleRequest struct {
	BytesPerLine int    `json:"bytesPerLine"`
	Text         string `json:"text"`
}

type TransformArticleResponse struct {
	Text string `json:"text"`
}

// TransformArticle godoc
// @Summary Transform Article
// @Description Transform Article, including split by bytes...etc
// @Tags ArticleProcessor
// @param input body TransformArticleRequest true "Transform data"
// @Accept json
// @Produce json
// @Success 200 {object} TransformArticleResponse
// @Failure 400 {object} problems.DefaultProblem
// @Router /article_process [post]
func (c *Controller) TransformArticle(ctx *gin.Context) {
	var req TransformArticleRequest
	err := ctx.Bind(&req)
	if err != nil {
		err400 := problems.NewDetailedProblem(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, err400)
		return
	}
	res, err := c.ArticleProcessorService.Transform(IrisAPIs.ProcessParameters{BytesPerLine: req.BytesPerLine}, req.Text)
	if err != nil {
		err500 := problems.NewDetailedProblem(http.StatusInternalServerError, err.Error())
		ctx.JSON(500, err500)
		return
	}
	ctx.JSON(200, TransformArticleResponse{
		Text: res,
	})
}
