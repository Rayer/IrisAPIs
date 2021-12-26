package IrisAPIsGRPC

import (
	"context"
	"github.com/Rayer/IrisAPIs"
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

func (a *ArticleProcessorServiceServiceGRPC) ProcessText(_ context.Context, req *ProcessTextRequest) (*ProcessTextResponse, error) {
	service := IrisAPIs.NewArticleProcessorContext()
	res, err := service.Transform(IrisAPIs.ProcessParameters{BytesPerLine: int(req.BytesPerLine)}, req.Text)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}
	return &ProcessTextResponse{
		ProcessedText: res,
	}, nil
}
