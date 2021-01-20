package IrisAPIsGRPC

import (
	"IrisAPIs"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ArticleProcessorServiceServiceGRPC struct {
	UnimplementedArticleProcessorServiceServer
	//service *IrisAPIs.ArticleProcessorService
}

func NewArticleProcessorServiceServiceGRPC() *ArticleProcessorServiceServiceGRPC {
	return &ArticleProcessorServiceServiceGRPC{}
}

func (a *ArticleProcessorServiceServiceGRPC) ProcessText(ctx context.Context, req *ProcessTextRequest) (*ProcessTextResponse, error) {
	service, err := IrisAPIs.NewArticleProcessorContext(IrisAPIs.ProcessParameters{BytesPerLine: int(req.BytesPerLine)})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}
	res, err := service.Transform(req.Text)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}

	return &ProcessTextResponse{
		ProcessedText: res,
	}, nil

}
