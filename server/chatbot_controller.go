package main

import (
	"github.com/gin-gonic/gin"
	"github.com/moogar0880/problems"
	"net/http"
)

type ChatbotReactResponse struct {
	Prompt          string   `json:"prompt"`
	Keywords        []string `json:"keywords"`
	InvalidKeywords []string `json:"invalid_keywords"`
	Message         string   `json:"message"`
	Error           string   `json:"error"`
	Next            string   `json:"next"`
}

type ChatbotResetUserResponse struct {
	User    string `json:"user"`
	Message string `json:"message"`
}

type ChatbotConversation struct {
	User  string `json:"user"`
	Input string `json:"input"`
}

// ChatBotReact godoc
// @Summary Reaction with ChatBot
// @Description Main Chatter interface for ChatBot
// @Tags ChatBot
// @Accept json
// @Produce json
// @param input body ChatbotConversation true "Input info"
// @Security ApiKeyAuth
// @Success 200 {object} ChatbotReactResponse
// @Failure 400 {object} problems.DefaultProblem
// @Router /chatbot [post]
func (c *Controller) ChatBotReact(ctx *gin.Context) {
	var conv ChatbotConversation
	err := ctx.BindJSON(&conv)

	utx, _ := c.ChatBotService.GetUserContext(conv.User)

	prompt, keywordsV, keywordsIv, err := utx.RenderMessageWithDetail()
	str, err := utx.HandleMessage(conv.Input)
	next, err := utx.RenderMessage()
	ctx.JSON(http.StatusOK, ChatbotReactResponse{
		Prompt:          prompt,
		Keywords:        keywordsV,
		InvalidKeywords: keywordsIv,
		Message:         str,
		Error:           err.Error(),
		Next:            next,
	})
}

// ChatBotResetUser godoc
// @Summary Reset user status to initial
// @Description Reset user status to initial
// @Tags ChatBot
// @Accept json
// @Produce json
// @Param user path string true "User name to reset"
// @Security ApiKeyAuth
// @Success 200 {object} ChatbotResetUserResponse
// @Failure 400 {object} problems.DefaultProblem
// @Router /chatbot/{user} [delete]
func (c *Controller) ChatBotResetUser(ctx *gin.Context) {
	user := ctx.Param("user")
	c.ChatBotService.ExpireUser(user, func() {
		ctx.JSON(http.StatusOK, ChatbotResetUserResponse{
			User:    user,
			Message: "ok",
		})
	}, func() {
		ctx.JSON(http.StatusBadRequest, problems.NewDetailedProblem(http.StatusBadRequest, "User "+user+" not found!"))
	})
}
