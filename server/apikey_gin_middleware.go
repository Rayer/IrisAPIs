package main

import (
	"IrisAPIs"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func ApiKeyCheckMiddleware(apiKeyService IrisAPIs.ApiKeyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.Query("apiKey")
		validKey := apiKeyService.ValidateApiKey(apiKey, IrisAPIs.QUERY_STRING)
		log.Debugf("Get request with ApiKey %s, which is %v\n", apiKey, validKey)
		c.Next()
	}
}
