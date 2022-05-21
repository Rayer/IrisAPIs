package main

import (
	"IrisAPIs"
	mock_IrisAPIs "IrisAPIs/server/mock"
	"encoding/json"
	"github.com/Pallinder/go-randomdata"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
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

func (a *ApiKeyControllerTestSuite) TestController_GetKey() {
	ctrl := gomock.NewController(a.T())
	defer ctrl.Finish()
	mockedModel := a.generateRandomApiKeyDataModel()
	mockService := mock_IrisAPIs.NewMockApiKeyService(ctrl)
	mockService.EXPECT().GetKeyModelById(gomock.Any(), *mockedModel.Id).Return(mockedModel, nil)
	mockService.EXPECT().GetKeyModelById(gomock.Any(), *mockedModel.Id+1).Return(nil, nil)
	mockService.EXPECT().GetKeyModelById(gomock.Any(), *mockedModel.Id+2).Return(nil, errors.Errorf("test error!"))

	type fields struct {
		ApiKeyService IrisAPIs.ApiKeyService
	}
	type args struct {
		Id string
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
				Id: strconv.Itoa(*mockedModel.Id),
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
				Id: strconv.Itoa(*mockedModel.Id + 1),
			},
			expectRetCode:   http.StatusNotFound,
			expectRetEntity: nil,
		},
		{
			name: "ShouldGetInternalError",
			fields: fields{
				ApiKeyService: mockService,
			},
			args: args{
				Id: strconv.Itoa(*mockedModel.Id + 2),
			},
			expectRetCode:   http.StatusInternalServerError,
			expectRetEntity: nil,
		},
		{
			name: "InvalidID",
			fields: fields{
				ApiKeyService: mockService,
			},
			args: args{
				Id: "abcd",
			},
			expectRetCode:   http.StatusBadRequest,
			expectRetEntity: nil,
		},
	}

	for _, tt := range tests {
		a.T().Run(tt.name, func(t *testing.T) {
			c := &Controller{
				ServiceMonolith: &IrisAPIs.ServiceMonolith{
					ApiKeyService: tt.fields.ApiKeyService,
				},
			}
			g, r := createGinTestItems()
			g.Request = httptest.NewRequest("GET", "/apiKey", nil)
			g.Params = []gin.Param{
				{
					Key:   "id",
					Value: tt.args.Id,
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
	mockService.EXPECT().IssueApiKey(gomock.Any(), "UTTestApp", true, true, "auto", false).Return("accc", nil)
	mockService.EXPECT().IssueApiKey(gomock.Any(), "TestMockErr", true, true, "auto", false).Return("", errors.New("unit test error"))
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
		{
			name: "ShouldGetInternalError",
			fields: fields{
				ApiKeyService: mockService,
			},
			args: args{
				payload: "{\"application\": \"TestMockErr\", \"use_in_header\": true,  \"use_in_query_param\": true}",
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		a.T().Run(tt.name, func(t *testing.T) {
			c := &Controller{
				ServiceMonolith: &IrisAPIs.ServiceMonolith{
					ApiKeyService: tt.fields.ApiKeyService,
				},
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
	mockService.EXPECT().GetAllKeys(gomock.Any()).Return([]*IrisAPIs.ApiKeyDataModel{
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
				ServiceMonolith: &IrisAPIs.ServiceMonolith{
					ApiKeyService: tt.fields.ApiKeyService,
				},
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

func (a *ApiKeyControllerTestSuite) TestController_GetApiUsage() {
	ctrl := gomock.NewController(a.T())
	defer ctrl.Finish()
	mockService := mock_IrisAPIs.NewMockApiKeyService(ctrl)

	type args struct {
		ctx *gin.Context
	}
	tests := []struct {
		name string
		args args
	}{
		{},
	}
	for _, tt := range tests {
		a.T().Run(tt.name, func(t *testing.T) {
			_ = &Controller{
				ServiceMonolith: &IrisAPIs.ServiceMonolith{
					ApiKeyService: mockService,
				},
			}
		})
	}
}
