//go:generate go get github.com/golang/mock/mockgen@v1.4.4
//go:generate ${GOPATH}/bin/mockgen -source ../ApiKeyContext.go -destination mock/ApiKeyContext.go
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
	"strings"
	"testing"
)

type ApiKeyControllerTestSuite struct {
	suite.Suite
	controller *Controller
}

func TestApiKeyControllerTestSuite(t *testing.T) {
	suite.Run(t, new(ApiKeyControllerTestSuite))
}

func (a *ApiKeyControllerTestSuite) SetupTest() {
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

	g, res := createGinTestItems()
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
	c.GetAllKeys(g)
	assert.Equal(a.T(), http.StatusOK, res.Code)
	var result []*ApiKeyBrief
	err := json.NewDecoder(res.Body).Decode(&result)
	assert.Equal(a.T(), nil, err)
	assert.Equal(a.T(), 2, len(result))
}

func (a *ApiKeyControllerTestSuite) TestController_GetKey() {
	ctrl := gomock.NewController(a.T())
	defer ctrl.Finish()
	mockedModel := a.generateRandomApiKeyDataModel()
	mockService := mock_IrisAPIs.NewMockApiKeyService(ctrl)
	mockService.EXPECT().GetKeyModelById(*mockedModel.Id).Return(mockedModel, nil)
	mockService.EXPECT().GetKeyModelById(*mockedModel.Id+1).Return(nil, nil)

	type fields struct {
		ApiKeyService IrisAPIs.ApiKeyService
	}
	type args struct {
		Id int
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		expectRetCode   int
		expectRetEntity *ApiKeyDetail
	}{
		{
			name: "ShouldFindEntry",
			fields: fields{
				ApiKeyService: mockService,
			},
			args: args{
				Id: *mockedModel.Id,
			},
			expectRetCode: http.StatusOK,
			expectRetEntity: &ApiKeyDetail{
				ApiKeyBrief: ApiKeyBrief{
					Id:         *mockedModel.Id,
					Key:        *mockedModel.Key,
					Privileged: *mockedModel.Privileged,
				},
				IssueBy:     "UTTester",
				Application: "UTTester",
			},
		},
		{
			name: "ShouldNotFindEntry",
			fields: fields{
				ApiKeyService: mockService,
			},
			args: args{
				Id: *mockedModel.Id + 1,
			},
			expectRetCode:   http.StatusNotFound,
			expectRetEntity: nil,
		},
	}

	for _, tt := range tests {
		a.T().Run(tt.name, func(t *testing.T) {
			c := &Controller{
				ApiKeyService: tt.fields.ApiKeyService,
			}
			g, r := createGinTestItems()
			id := strconv.Itoa(tt.args.Id)
			g.Params = []gin.Param{
				{
					Key:   "id",
					Value: id,
				},
			}
			c.GetKey(g)
			assert.Equal(a.T(), tt.expectRetCode, r.Code)

			if tt.expectRetEntity != nil {
				var result ApiKeyDetail
				err := json.NewDecoder(r.Body).Decode(&result)
				if err != nil {
					a.T().Fail()
				}
				assert.Equal(a.T(), *mockedModel.Key, result.Key)
				assert.Equal(a.T(), *mockedModel.Id, result.Id)
				assert.Equal(a.T(), *mockedModel.Privileged, result.Privileged)
			}
		})

	}
}

func (a *ApiKeyControllerTestSuite) TestController_IssueApiKey() {
	ctrl := gomock.NewController(a.T())
	defer ctrl.Finish()
	mockService := mock_IrisAPIs.NewMockApiKeyService(ctrl)
	mockService.EXPECT().IssueApiKey("UTTestApp", true, true, "auto", false).Return("accc", nil)

	type fields struct {
		ApiKeyService IrisAPIs.ApiKeyService
	}
	type args struct {
		payload string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		expectedStatus int
	}{
		{
			name: "SuccessIssue",
			fields: fields{
				ApiKeyService: mockService,
			},
			args: args{
				payload: "{\"application\": \"UTTestApp\", \"use_in_header\": true,  \"use_in_query_param\": true}",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "FailedBadPayload",
			fields: fields{
				ApiKeyService: mockService,
			},
			args: args{
				payload: "Ohmygodhow!",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		a.T().Run(tt.name, func(t *testing.T) {
			c := &Controller{
				ApiKeyService: tt.fields.ApiKeyService,
			}
			g, r := createGinTestItems()
			payload := strings.NewReader(tt.args.payload)
			g.Request = httptest.NewRequest("POST", "/apiKey", payload)
			c.IssueApiKey(g)
			assert.Equal(t, tt.expectedStatus, r.Code)
		})
	}
}

func (a *ApiKeyControllerTestSuite) TestController_GetAllKeys() {
	ctrl := gomock.NewController(a.T())
	defer ctrl.Finish()
	mockService := mock_IrisAPIs.NewMockApiKeyService(ctrl)
	mockService.EXPECT().GetAllKeys().Return([]*IrisAPIs.ApiKeyDataModel{
		a.generateRandomApiKeyDataModel(),
		a.generateRandomApiKeyDataModel(),
		a.generateRandomApiKeyDataModel(),
	}, nil)

	type fields struct {
		ApiKeyService IrisAPIs.ApiKeyService
	}
	tests := []struct {
		name             string
		fields           fields
		expectedStatus   int
		expectedArrLenth int
	}{
		{
			name: "NormalTest",
			fields: fields{
				ApiKeyService: mockService,
			},
			expectedStatus:   http.StatusOK,
			expectedArrLenth: 3,
		},
	}
	for _, tt := range tests {
		a.T().Run(tt.name, func(t *testing.T) {
			c := &Controller{
				ApiKeyService: tt.fields.ApiKeyService,
			}
			g, r := createGinTestItems()
			c.GetAllKeys(g)
			assert.Equal(t, tt.expectedStatus, r.Code)
			var decoded []ApiKeyBrief
			err := json.NewDecoder(r.Body).Decode(&decoded)
			if err != nil {
				t.Fail()
			}
			assert.Equal(t, tt.expectedArrLenth, len(decoded))
		})
	}
}
