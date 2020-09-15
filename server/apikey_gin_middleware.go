package main

import (
	"IrisAPIs"
	"github.com/gin-gonic/gin"
	"github.com/moogar0880/problems"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func ApiKeyCheckMiddleware(apiKeyService IrisAPIs.ApiKeyService, enforceCheckApiKey bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var keyLocation IrisAPIs.ApiKeyLocation
		keyLocation = IrisAPIs.QueryString
		apiKey := c.Query("apiKey")
		if apiKey == "" {
			apiKey = c.GetHeader("apiKey")
			keyLocation = IrisAPIs.Header
		}
		validKey := apiKeyService.ValidateApiKey(apiKey, keyLocation)
		if !validKey && enforceCheckApiKey {
			c.JSON(http.StatusUnauthorized, problems.NewDetailedProblem(http.StatusUnauthorized, "Not authorize for this resource"))
			c.Abort()
			return
		}
		log.Debugf("Get request with ApiKey %s, which is %v", apiKey, validKey)
		c.Next()
	}
}
