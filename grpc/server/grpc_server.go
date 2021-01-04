package main

import (
	"IrisAPIs"
	IrisAPIsGRPC "IrisAPIs/grpc"
)

func main() {
	conf := IrisAPIs.NewConfiguration()
	s := new(IrisAPIsGRPC.GRPCServerRoutine)
	s.Run(conf)
}
