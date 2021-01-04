package IrisAPIsGRPC

import (
	"IrisAPIs"
	"context"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

type Events struct {
}

type ServerRoutine interface {
	Run()
	RunDetach() chan Events
}

type GRPCServerRoutine struct {
}

func (g *GRPCServerRoutine) Run(conf *IrisAPIs.Configuration) {
	g.runImpl(conf)
}

func (g *GRPCServerRoutine) RunDetach(ctx context.Context, conf *IrisAPIs.Configuration) chan Events {
	events := make(chan Events)
	go func() {
		g.runImpl(conf)
	}()
	return events
}

func (g *GRPCServerRoutine) runImpl(conf *IrisAPIs.Configuration) {
	lis, err := net.Listen("tcp", conf.GRPCServerHost)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	RegisterApiKeyServiceServer(s, NewApiKeyServiceGRPC(conf.ConnectionString))
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}