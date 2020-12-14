//go:generate ${GOPATH}/bin/mockgen -source ../CurrencyContext.go -destination mock/CurrencyContext.go
package main

import (
	"IrisAPIs"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"net/http/httptest"
	"testing"
)

type TestCurrencyControllerSuite struct {
	suite.Suite
	testGin  *gin.Context
	recorder *httptest.ResponseRecorder
}

func TestController_ConvertCurrency(t *testing.T) {
	type fields struct {
		SystemDefaultController SystemDefaultController
		ChatBotService          *IrisAPIs.ChatbotContext
		CurrencyService         *IrisAPIs.CurrencyContext
		DatabaseContext         *IrisAPIs.DatabaseContext
		IpNationService         *IrisAPIs.IpNationContext
		ApiKeyService           IrisAPIs.ApiKeyService
		ServiceMgmt             IrisAPIs.ServiceManagement
	}
	type args struct {
		ctx *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = &Controller{
				SystemDefaultController: tt.fields.SystemDefaultController,
				ChatBotService:          tt.fields.ChatBotService,
				CurrencyService:         tt.fields.CurrencyService,
				DatabaseContext:         tt.fields.DatabaseContext,
				IpNationService:         tt.fields.IpNationService,
				ApiKeyService:           tt.fields.ApiKeyService,
				ServiceMgmt:             tt.fields.ServiceMgmt,
			}
		})
	}
}

func TestController_GetCurrencyRaw(t *testing.T) {
	type fields struct {
		SystemDefaultController SystemDefaultController
		ChatBotService          *IrisAPIs.ChatbotContext
		CurrencyService         *IrisAPIs.CurrencyContext
		DatabaseContext         *IrisAPIs.DatabaseContext
		IpNationService         *IrisAPIs.IpNationContext
		ApiKeyService           IrisAPIs.ApiKeyService
		ServiceMgmt             IrisAPIs.ServiceManagement
	}
	type args struct {
		ctx *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = &Controller{
				SystemDefaultController: tt.fields.SystemDefaultController,
				ChatBotService:          tt.fields.ChatBotService,
				CurrencyService:         tt.fields.CurrencyService,
				DatabaseContext:         tt.fields.DatabaseContext,
				IpNationService:         tt.fields.IpNationService,
				ApiKeyService:           tt.fields.ApiKeyService,
				ServiceMgmt:             tt.fields.ServiceMgmt,
			}
		})
	}
}

func TestController_SyncData(t *testing.T) {
	type fields struct {
		SystemDefaultController SystemDefaultController
		ChatBotService          *IrisAPIs.ChatbotContext
		CurrencyService         *IrisAPIs.CurrencyContext
		DatabaseContext         *IrisAPIs.DatabaseContext
		IpNationService         *IrisAPIs.IpNationContext
		ApiKeyService           IrisAPIs.ApiKeyService
		ServiceMgmt             IrisAPIs.ServiceManagement
	}
	type args struct {
		ctx *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = &Controller{
				SystemDefaultController: tt.fields.SystemDefaultController,
				ChatBotService:          tt.fields.ChatBotService,
				CurrencyService:         tt.fields.CurrencyService,
				DatabaseContext:         tt.fields.DatabaseContext,
				IpNationService:         tt.fields.IpNationService,
				ApiKeyService:           tt.fields.ApiKeyService,
				ServiceMgmt:             tt.fields.ServiceMgmt,
			}
		})
	}
}
