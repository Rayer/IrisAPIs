package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"IrisAPIs"
	"IrisAPIs/server/graph/generated"
	"IrisAPIs/server/graph/model"
	"context"
	"fmt"
	"github.com/docker/distribution/uuid"
)

func (r *mutationResolver) PostArticleProcess(ctx context.Context, mainTransformArticleRequestInput model.MainTransformArticleRequestInput) (*model.MainTransformArticleResponse, error) {
	return nil, fmt.Errorf("not supported")
}

func (r *mutationResolver) PostCurrencyConvert(ctx context.Context, mainCurrencyConvertInput model.MainCurrencyConvertInput) (*model.MainCurrencyConvert, error) {
	return nil, fmt.Errorf("not supported")
}

func (r *mutationResolver) MutationViewerAPIKey(ctx context.Context, apiKey string) (*model.MutationViewerAPIKey, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) MutationViewerAnyAuth(ctx context.Context, apiKeyAuth *model.APIKeyAuthInput) (*model.MutationViewerAnyAuth, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) MainGetServiceStatusByIDResponse(ctx context.Context, id string) (*model.MainGetServiceStatusByIDResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) MainIPNationMyIPResponse(ctx context.Context) (*model.MainIPNationMyIPResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) MainPingResponse(ctx context.Context) (*model.MainPingResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) PbsRecent(ctx context.Context, format *string, period *string) ([]*model.MainGetRecentPBSDataResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Service(ctx context.Context) ([]*model.MainGetServiceStatusByIDResponse, error) {
	var ret []*model.MainGetServiceStatusByIDResponse
	serviceStatusList := r.Services.ServiceMgmt.CheckAllServerStatus(ctx)
	for _, s := range serviceStatusList {
		ret = append(ret, &model.MainGetServiceStatusByIDResponse{
			ID:      IrisAPIs.PValue(s.ID.String()),
			Message: IrisAPIs.PValue(s.Message),
			Name:    IrisAPIs.PValue(s.Name),
			Status:  IrisAPIs.PValue(string(s.Status)),
			Type:    IrisAPIs.PValue(s.ServiceType),
		})
	}
	return ret, nil
}

func (r *queryResolver) ServiceLogs(ctx context.Context, id string) (*string, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	logs, err := r.Services.ServiceMgmt.GetLogs(ctx, uid)
	if err != nil {
		return nil, err
	}
	return &logs, nil
}

func (r *queryResolver) ViewerAPIKey(ctx context.Context, apiKey string) (*model.ViewerAPIKey, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ViewerAnyAuth(ctx context.Context, apiKeyAuth *model.APIKeyAuthInput) (*model.ViewerAnyAuth, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
