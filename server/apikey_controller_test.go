package main

import (
	"IrisAPIs"
	mock_IrisAPIs "IrisAPIsServer/mock"
	"encoding/json"
	"github.com/Pallinder/go-randomdata"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type ApiKeyControllerTestSuite struct {
	suite.Suite
	gin              *gin.Context
	controller       *Controller
	responseRecorder *httptest.ResponseRecorder
}

func TestApiKeyControllerTestSuite(t *testing.T) {
	suite.Run(t, new(ApiKeyControllerTestSuite))
}

func (a *ApiKeyControllerTestSuite) SetupTest() {
	a.responseRecorder = httptest.NewRecorder()
	a.gin, _ = gin.CreateTestContext(a.responseRecorder)
	a.controller = &Controller{}
}

func (a *ApiKeyControllerTestSuite) generateRandomApiKeyDataModel() *IrisAPIs.ApiKeyDataModel {
	return &IrisAPIs.ApiKeyDataModel{
		Id:          IrisAPIs.PInt(randomdata.Number(0, 100)),
		Key:         IrisAPIs.PString(randomdata.Alphanumeric(16)),
		UseInHeader: IrisAPIs.PBool(randomdata.Boolean()),
		UseInQuery:  IrisAPIs.PBool(randomdata.Boolean()),
		Application: IrisAPIs.PString("IrisAPIsServerUT"),
		Issuer:      IrisAPIs.PString(randomdata.SillyName()),
		IssueDate:   nil,
		Privileged:  IrisAPIs.PBool(randomdata.Boolean()),
		Expiration:  nil,
	}

}

func (a *ApiKeyControllerTestSuite) TestController_GetApiUsage() {

	ctrl := gomock.NewController(a.T())
	defer ctrl.Finish()
	mockService := mock_IrisAPIs.NewMockApiKeyService(ctrl)
	c := &Controller{
		ApiKeyService: mockService,
	}
	mockService.EXPECT().GetAllKeys().Return([]*IrisAPIs.ApiKeyDataModel{
		a.generateRandomApiKeyDataModel(),
		a.generateRandomApiKeyDataModel(),
	}, nil)
	c.GetAllKeys(a.gin)
	assert.Equal(a.T(), http.StatusOK, a.responseRecorder.Code)
	var result []*ApiKeyBrief
	err := json.NewDecoder(a.responseRecorder.Body).Decode(&result)
	assert.Equal(a.T(), nil, err)
	assert.Equal(a.T(), 2, len(result))
}

func (a *ApiKeyControllerTestSuite) TestController_GetKey() {
	ctrl := gomock.NewController(a.T())
	defer ctrl.Finish()
	mockService := mock_IrisAPIs.NewMockApiKeyService(ctrl)
	c := &Controller{
		ApiKeyService: mockService,
	}

	mockedModel := a.generateRandomApiKeyDataModel()
	mockService.EXPECT().GetKeyModelById(*mockedModel.Id).Return(mockedModel, nil)
	mockService.EXPECT().GetKeyModelById(*mockedModel.Id+1).Return(nil, nil)
	a.gin.Params = []gin.Param{
		{
			Key:   "id",
			Value: strconv.Itoa(*mockedModel.Id),
		},
	}
	c.GetKey(a.gin)
	assert.Equal(a.T(), http.StatusOK, a.responseRecorder.Code)
	a.responseRecorder.Flush()

	a.gin.Params = []gin.Param{
		{
			Key:   "id",
			Value: strconv.Itoa(*mockedModel.Id + 1),
		},
	}
	c.GetKey(a.gin)
	assert.Equal(a.T(), http.StatusNotFound, a.responseRecorder.Code)

}
