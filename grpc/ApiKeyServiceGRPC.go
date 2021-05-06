package IrisAPIsGRPC

import (
	"IrisAPIs"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type ApiKeyServiceGRPC struct {
	UnimplementedApiKeyServiceServer
	service IrisAPIs.ApiKeyService
}

func NewApiKeyServiceGRPC(connectionString string) *ApiKeyServiceGRPC {
	return &ApiKeyServiceGRPC{
		service: IrisAPIs.NewApiKeyService(func() *IrisAPIs.DatabaseContext {
			ret, _ := IrisAPIs.NewDatabaseContext(connectionString, false)
			return ret
		}()),
	}
}

func (a *ApiKeyServiceGRPC) IssueApiKey(ctx context.Context, r *IssueApiKeyRequest) (*IssueApiKeyResponse, error) {
	ret, err := a.service.IssueApiKey(context.TODO(), r.Application, r.UseInHandler, r.UseInQuery, r.Issuer, r.Privileged)
	return &IssueApiKeyResponse{
		ApiKey: ret,
	}, status.Errorf(codes.Internal, "%s", err.Error())
}
func (a *ApiKeyServiceGRPC) ValidateApiKey(ctx context.Context, r *ValidateApiKeyRequest) (*ValidateApiKeyResponse, error) {
	_, previlegeLevel := a.service.ValidateApiKey(context.TODO(), r.Key, IrisAPIs.ApiKeyLocation(r.ApiKeyLocation))
	return &ValidateApiKeyResponse{
		//2 is for offset between PrivilegeLevel(protobuf) and ApiKeyPrivilegeLevel
		PrivilegeLevel: PrivilegeLevel(previlegeLevel + 2),
	}, nil
}

func ApiKeyDataModelToGRPC(v *IrisAPIs.ApiKeyDataModel) *ApiKeyDetail {
	var expiration *wrapperspb.Int64Value
	if v.Expiration == nil {
		expiration = nil
	} else {
		expiration = &wrapperspb.Int64Value{
			Value: v.Expiration.Unix(),
		}
	}

	return &ApiKeyDetail{
		Id:          int64(*v.Id),
		Key:         *v.Key,
		UseInHeader: *v.UseInQuery,
		UseInQuery:  *v.UseInHeader,
		Application: *v.Application,
		Issuer:      *v.Issuer,
		IssueDate:   v.IssueDate.Unix(),
		Privileged:  *v.Privileged,
		Expiration:  expiration,
	}
}

func (a *ApiKeyServiceGRPC) GetAllKeys(ctx context.Context, r *GetAllKeysRequest) (*GetAllKeysResponse, error) {
	keys, err := a.service.GetAllKeys(context.TODO())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}

	entries := make([]*ApiKeyDetail, 0)

	for _, v := range keys {
		entries = append(entries, ApiKeyDataModelToGRPC(v))
	}

	return &GetAllKeysResponse{
		Entries: entries,
	}, nil
}

func (a *ApiKeyServiceGRPC) GetKeyById(ctx context.Context, r *GetKeyByIdRequest) (*GetKeyByIdResponse, error) {
	id := r.GetId()
	e, err := a.service.GetKeyModelById(context.TODO(), int(id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}
	return &GetKeyByIdResponse{
		Entry: ApiKeyDataModelToGRPC(e),
	}, nil
}

func (a *ApiKeyServiceGRPC) SetExpired(ctx context.Context, r *SetExpiredRequest) (*SetExpiredResponse, error) {
	err := a.service.SetExpire(context.TODO(), int(r.Id), r.IsExpired)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}
	return &SetExpiredResponse{}, nil
}
