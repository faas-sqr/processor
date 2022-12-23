package grpc_client

import (
	"google.golang.org/grpc"
	"log"
	"processor/pkg/grpc-client/proto/function"
)

var ClientMap map[string]function.HandleEventClient

func init() {
	ClientMap = make(map[string]function.HandleEventClient)
}

func ClientInit(functionName, port string) {
	conn, err := grpc.Dial(":"+port, grpc.WithInsecure())
	if err != nil {
		//重试机制
		log.Fatalf("did not connect: %v", err)
	}
	//defer conn.Close()

	c := function.NewHandleEventClient(conn)
	ClientMap[functionName] = c
}
