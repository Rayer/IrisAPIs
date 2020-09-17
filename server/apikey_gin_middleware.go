package main

import (
	"IrisAPIs"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/moogar0880/problems"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type ApiKeyValidator interface {
	GetMiddleware() gin.HandlerFunc
	//RegisterPrivilegeRoute(group string, privilegeLevel IrisAPIs.ApiKeyPrivilegeLevel, entry func(string, ...gin.HandlerFunc) gin.IRoutes, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
	RegisterPrivilegeRoute(fullPath string, method string, level IrisAPIs.ApiKeyPrivilegeLevel)
	FetchPrivilegeLevel(fullPath string, method string) IrisAPIs.ApiKeyPrivilegeLevel
}

type ApiKeyValidatorContext struct {
	privilegeRoutes map[string]IrisAPIs.ApiKeyPrivilegeLevel
	apiKeyService   IrisAPIs.ApiKeyService
	EnforceApiKey   bool
}

func NewApiKeyValidator(apiKeyService IrisAPIs.ApiKeyService, enforceCheckApiKey bool) ApiKeyValidator {
	return &ApiKeyValidatorContext{apiKeyService: apiKeyService, EnforceApiKey: enforceCheckApiKey, privilegeRoutes: make(map[string]IrisAPIs.ApiKeyPrivilegeLevel)}
}

func (a *ApiKeyValidatorContext) GetMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var keyLocation IrisAPIs.ApiKeyLocation
		keyLocation = IrisAPIs.QueryString
		apiKey := c.Query("apiKey")
		if apiKey == "" {
			apiKey = c.GetHeader("apiKey")
			keyLocation = IrisAPIs.Header
		}

		path := c.FullPath()
		method := c.Request.Method

		pathPrivilege := a.FetchPrivilegeLevel(path, method)

		fmt.Printf("Privilege Map : %+v\n", a.privilegeRoutes)
		fmt.Printf("Path : %s, Path Privilege : %d\n", path, pathPrivilege)

		validKey := a.apiKeyService.ValidateApiKey(apiKey, keyLocation)

		if validKey < pathPrivilege && a.EnforceApiKey {
			c.JSON(http.StatusUnauthorized, problems.NewDetailedProblem(http.StatusUnauthorized, "Not authorize for this resource"))
			c.Abort()
			return
		}
		log.Debugf("Get request with ApiKey %s, which is %v", apiKey, validKey)
		c.Next()
	}
}

func (a *ApiKeyValidatorContext) RegisterPrivilegeRoute(fullPath string, method string, level IrisAPIs.ApiKeyPrivilegeLevel) {
	storeKey := fullPath + "+" + method
	a.privilegeRoutes[storeKey] = level
}

func (a *ApiKeyValidatorContext) FetchPrivilegeLevel(fullPath string, method string) IrisAPIs.ApiKeyPrivilegeLevel {
	storeKey := fullPath + "+" + method
	if p, ok := a.privilegeRoutes[storeKey]; ok {
		return p
	} else {
		return IrisAPIs.ApiKeyNotPresented
	}
}

type AKWrappedEngine struct {
	*gin.Engine
	validator ApiKeyValidator
}

func NewAKWrappedEngine(engine *gin.Engine, validator ApiKeyValidator) *AKWrappedEngine {
	return &AKWrappedEngine{Engine: engine, validator: validator}
}

func (e *AKWrappedEngine) Group(relativePath string, handlers ...gin.HandlerFunc) *AKGroup {
	return NewAKGroup(relativePath, e.Engine, e.validator, handlers...)
}

type AKGroup struct {
	relativePath string
	wrapped      *gin.RouterGroup
	validator    ApiKeyValidator
}

func NewAKGroup(relativePath string, engine *gin.Engine, validator ApiKeyValidator, handlers ...gin.HandlerFunc) *AKGroup {
	return &AKGroup{relativePath: relativePath, wrapped: engine.Group(relativePath, handlers...), validator: validator}
}

func (group *AKGroup) POST(relativePath string, level IrisAPIs.ApiKeyPrivilegeLevel, handlers ...gin.HandlerFunc) gin.IRoutes {
	group.registerPrivilegeEndpoint(group.relativePath+relativePath, http.MethodPost, level)
	return group.wrapped.POST(relativePath, handlers...)
}

func (group *AKGroup) GET(relativePath string, level IrisAPIs.ApiKeyPrivilegeLevel, handlers ...gin.HandlerFunc) gin.IRoutes {
	group.registerPrivilegeEndpoint(group.relativePath+relativePath, http.MethodGet, level)
	return group.wrapped.GET(relativePath, handlers...)
}

func (group *AKGroup) PUT(relativePath string, level IrisAPIs.ApiKeyPrivilegeLevel, handlers ...gin.HandlerFunc) gin.IRoutes {
	group.registerPrivilegeEndpoint(group.relativePath+relativePath, http.MethodPut, level)
	return group.wrapped.PUT(relativePath, handlers...)
}

func (group *AKGroup) DELETE(relativePath string, level IrisAPIs.ApiKeyPrivilegeLevel, handlers ...gin.HandlerFunc) gin.IRoutes {
	group.registerPrivilegeEndpoint(group.relativePath+relativePath, http.MethodDelete, level)
	return group.wrapped.DELETE(relativePath, handlers...)
}

func (group *AKGroup) registerPrivilegeEndpoint(fullPath string, httpMethod string, level IrisAPIs.ApiKeyPrivilegeLevel) {
	group.validator.RegisterPrivilegeRoute(fullPath, httpMethod, level)
}