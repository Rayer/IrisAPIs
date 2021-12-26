package cmd

import (
	"context"
	"github.com/Rayer/IrisAPIs"
	IrisAPIsGRPC "github.com/Rayer/IrisAPIs/grpc"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"time"
)

type LocalDataSource struct {
	IrisAPIs.ApiKeyService
}

type GRPCDataSource struct {
	client IrisAPIsGRPC.ApiKeyServiceClient
	ctx    context.Context
}

func NewGRPCDataSource(grpcServer string) IrisAPIs.ApiKeyService {
	conn, err := grpc.Dial(grpcServer, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Connection fail：%v", err)
	}
	ret := new(GRPCDataSource)
	ret.client = IrisAPIsGRPC.NewApiKeyServiceClient(conn)
	ret.ctx = context.Background()
	return ret
}

func (G *GRPCDataSource) IssueApiKey(ctx context.Context, application string, useInHeader bool, useInQuery bool, issuer string, privileged bool) (string, error) {
	ret, err := G.client.IssueApiKey(G.ctx, &IrisAPIsGRPC.IssueApiKeyRequest{
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

func (G *GRPCDataSource) ValidateApiKey(ctx context.Context, key string, embeddedIn IrisAPIs.ApiKeyLocation) (int, IrisAPIs.ApiKeyPrivilegeLevel) {
	ret, err := G.client.ValidateApiKey(G.ctx, &IrisAPIsGRPC.ValidateApiKeyRequest{
		Key:            key,
		ApiKeyLocation: int64(embeddedIn),
	})
	if err != nil {
		return -1, IrisAPIs.ApiKeyPrivilegeLevel(IrisAPIs.ApiKeyNotPresented)
	}

	return -1, IrisAPIs.ApiKeyPrivilegeLevel(ret.PrivilegeLevel)
}

func (G *GRPCDataSource) RecordActivity(ctx context.Context, path string, method string, key string, location IrisAPIs.ApiKeyLocation, ip string) {
	panic("implement me")
}

func (G *GRPCDataSource) GetAllKeys(context.Context) ([]*IrisAPIs.ApiKeyDataModel, error) {
	res, err := G.client.GetAllKeys(G.ctx, &IrisAPIsGRPC.GetAllKeysRequest{})
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

func (G *GRPCDataSource) GetKeyModelById(ctx context.Context, id int) (*IrisAPIs.ApiKeyDataModel, error) {
	panic("implement me")
}

func (G *GRPCDataSource) SetExpire(ctx context.Context, keyId int, expire bool) error {
	_, err := G.client.SetExpired(G.ctx, &IrisAPIsGRPC.SetExpiredRequest{
		Id:        int64(keyId),
		IsExpired: expire,
	})
	return err
}

func (G *GRPCDataSource) GetKeyUsageById(ctx context.Context, id int, from *time.Time, to *time.Time) ([]*IrisAPIs.ApiKeyAccess, error) {
	resp, err := G.client.GetKeyUsageById(G.ctx, &IrisAPIsGRPC.GetKeyUsageByIdRequest{
		Id:   int64(id),
		From: IrisAPIs.PGTimestamp(from),
		To:   IrisAPIs.PGTimestamp(to),
	})
	if err != nil {
		return nil, err
	}
	ret := make([]*IrisAPIs.ApiKeyAccess, 0)
	for _, v := range resp.Entries {
		ret = append(ret, &IrisAPIs.ApiKeyAccess{
			Id:        IrisAPIs.PInt(int(v.Id)),
			ApiKeyRef: IrisAPIs.PInt(int(v.ApiKeyRef)),
			Fullpath:  IrisAPIs.PString(v.FullPath),
			Method:    IrisAPIs.PString(v.Method),
			Ip:        IrisAPIs.PString(v.Ip),
			Nation:    IrisAPIs.PString(v.Nation),
			Timestamp: IrisAPIs.PTime(time.Unix(v.Timestamp, 0)),
		})
	}
	return ret, nil
}

func (G *GRPCDataSource) GetKeyUsageByPath(ctx context.Context, path string, exactMatch bool, from *time.Time, to *time.Time) ([]*IrisAPIs.ApiKeyAccess, error) {

	resp, err := G.client.GetKeyUsageByPath(G.ctx, &IrisAPIsGRPC.GetKeyUsageByPathRequest{
		Path:       path,
		ExactMatch: exactMatch,
		From:       IrisAPIs.PGTimestamp(from),
		To:         IrisAPIs.PGTimestamp(to),
	})
	if err != nil {
		return nil, err
	}
	ret := make([]*IrisAPIs.ApiKeyAccess, 0)
	for _, v := range resp.Entries {
		ret = append(ret, &IrisAPIs.ApiKeyAccess{
			Id:        IrisAPIs.PInt(int(v.Id)),
			ApiKeyRef: IrisAPIs.PInt(int(v.ApiKeyRef)),
			Fullpath:  IrisAPIs.PString(v.FullPath),
			Method:    IrisAPIs.PString(v.Method),
			Ip:        IrisAPIs.PString(v.Ip),
			Nation:    IrisAPIs.PString(v.Nation),
			Timestamp: IrisAPIs.PTime(time.Unix(v.Timestamp, 0)),
		})
	}
	return ret, nil
}
