package main

import (
	"github.com/Rayer/IrisAPIs"
	IrisAPIsGRPC "github.com/Rayer/IrisAPIs/grpc"
)

func main() {
	conf := IrisAPIs.NewConfiguration()
	s := new(IrisAPIsGRPC.GRPCServerRoutine)
	s.Run(conf)
}
