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
	"strconv"
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
