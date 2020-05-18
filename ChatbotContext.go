package IrisAPIs

import (
	"fmt"
	ChatBot "github.com/rayer/chatbot"
)

type ChatbotContext struct {
	conf   ChatBot.Configuration
	ctxMgr *ChatBot.ContextManager
}

type ChatbotConversion struct {
	User  string `json:"user"`
	Input string `json:"input"`
}

func NewChatbotContext() *ChatbotContext {
	ret := ChatbotContext{}
	ret.conf = ChatBot.Configuration{
		ResetTimerSec: 300,
		KeywordFormatter: func(fullMessage string, keyword string, isValidKeyword bool) string {
			return fmt.Sprintf("[%s]", keyword)
		},
	}
	ret.ctxMgr = ChatBot.NewContextManagerWithConfig(&ret.conf)
	return &ret
}

func (c *ChatbotContext) GetUserContext(user string) (userContext *ChatBot.UserContext, existing bool) {
	var ret *ChatBot.UserContext
	existing = true
	ret = c.ctxMgr.GetUserContext(user)
	if ret == nil {
		existing = false
		ret = c.ctxMgr.CreateUserContext(user, func() ChatBot.Scenario {
			return &RootScenario{}
		})
	}

	return ret, existing
}

func (c *ChatbotContext) ExpireUser(user string, foundUser func(), notfound func()) {
	ret := c.ctxMgr.GetUserContext(user)
	if ret == nil {
		if notfound != nil {
			notfound()
		}
	} else {
		c.ctxMgr.ExpireUser(user)
		if foundUser != nil {
			foundUser()
		}
	}
}
