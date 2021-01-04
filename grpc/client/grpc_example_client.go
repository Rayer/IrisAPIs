package main

import (
	"IrisAPIs"
	IrisAPIsGRPC "IrisAPIs/grpc"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	conf := IrisAPIs.NewConfiguration()
	conn, err := grpc.Dial(conf.GRPCServerHost, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Connection fail：%v", err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Printf("Error closing connection : %s", err.Error())
		}
	}()

	c := IrisAPIsGRPC.NewApiKeyServiceClient(conn)
	ctx := context.Background()
	r, err := c.GetAllKeys(ctx, &IrisAPIsGRPC.GetAllKeysRequest{})
	keys := make([]string, 0)
	for i, v := range r.Entries {
		keys = append(keys, v.Key)
		if i > 10 {
			break
		}
	}

	for _, key := range keys {
		r, _ := c.ValidateApiKey(ctx, &IrisAPIsGRPC.ValidateApiKeyRequest{
			Key:            key,
			ApiKeyLocation: 1,
		})
		fmt.Println(r.PrivilegeLevel.Enum())
	}

	fmt.Printf("%+v", r.Entries)
}
