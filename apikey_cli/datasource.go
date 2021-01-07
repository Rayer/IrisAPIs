package main

import (
	"IrisAPIs"
	IrisAPIsGRPC "IrisAPIs/grpc"
	"context"
	"time"
)

type LocalDataSource struct {
	IrisAPIs.ApiKeyService
}

type GRPCDataSource struct {
	provider *IrisAPIsGRPC.ApiKeyServiceGRPC
}

func (G *GRPCDataSource) IssueApiKey(application string, useInHeader bool, useInQuery bool, issuer string, privileged bool) (string, error) {
	ret, err := G.provider.IssueApiKey(context.TODO(), &IrisAPIsGRPC.IssueApiKeyRequest{
		Application:  application,
		UseInHandler: useInHeader,
		UseInQuery:   useInQuery,
		Issuer:       issuer,
		Privileged:   privileged,
	})

	if err != nil {
		return "", err
	}

	return ret.ApiKey, err
}

func (G *GRPCDataSource) ValidateApiKey(key string, embeddedIn IrisAPIs.ApiKeyLocation) IrisAPIs.ApiKeyPrivilegeLevel {
	ret, err := G.provider.ValidateApiKey(context.TODO(), &IrisAPIsGRPC.ValidateApiKeyRequest{
		Key:            key,
		ApiKeyLocation: int32(embeddedIn),
	})
	if err != nil {
		return IrisAPIs.ApiKeyPrivilegeLevel(IrisAPIs.ApiKeyNotPresented)
	}

	return IrisAPIs.ApiKeyPrivilegeLevel(ret.PrivilegeLevel)
}

func (G *GRPCDataSource) RecordActivity(path string, method string, key string, location IrisAPIs.ApiKeyLocation, ip string) {
	panic("implement me")
}

func (G *GRPCDataSource) GetAllKeys() ([]*IrisAPIs.ApiKeyDataModel, error) {
	res, err := G.provider.GetAllKeys(context.TODO(), &IrisAPIsGRPC.GetAllKeysRequest{})
	if err != nil {
		return nil, err
	}
	ret := make([]*IrisAPIs.ApiKeyDataModel, 0)
	for _, v := range res.Entries {
		ret = append(ret, &IrisAPIs.ApiKeyDataModel{
			Id:          IrisAPIs.PInt(int(v.Id)),
			Key:         &v.Key,
			UseInHeader: &v.UseInHeader,
			UseInQuery:  &v.UseInQuery,
			Application: &v.Application,
			Issuer:      &v.Issuer,
			IssueDate:   IrisAPIs.PTime(time.Unix(v.IssueDate, 0)),
			Privileged:  &v.Privileged,
			Expiration: func() *time.Time {
				if v.Expiration == nil {
					return nil
				}
				return IrisAPIs.PTime(time.Unix(v.Expiration.GetValue(), 0))
			}(),
		})
	}
	return ret, nil
}

func (G *GRPCDataSource) GetKeyModelById(id int) (*IrisAPIs.ApiKeyDataModel, error) {
	panic("implement me")
}

func (G *GRPCDataSource) SetExpire(keyId int, expire bool) error {
	panic("implement me")
}

func (G *GRPCDataSource) GetKeyUsageById(id int, from *time.Time, to *time.Time) ([]*IrisAPIs.ApiKeyAccess, error) {
	panic("implement me")
}

func (G *GRPCDataSource) GetKeyUsageByPath(path string, exactMatch bool, from *time.Time, to *time.Time) ([]*IrisAPIs.ApiKeyAccess, error) {
	panic("implement me")
}
