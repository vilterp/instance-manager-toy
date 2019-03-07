package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cockroachlabs/instance_manager/pure_manager/proto"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("hello this is the client")

	addr := "0.0.0.0:8888"

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("failed to dial:", err)
	}
	client := proto.NewGroupManagerClient(conn)

	ctx := context.Background()
	fmt.Println("get current spec:")
	resp, err := client.GetCurrentSpec(ctx, &proto.GetCurrentSpecRequest{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)

	fmt.Println("update:")
	resp2, err2 := client.UpdateSpec(ctx, &proto.UpdateSpecRequest{
		Spec: &proto.GroupSpec{
			NumInstances: 3,
		},
	})
	if err2 != nil {
		log.Fatal(err2)
	}
	fmt.Println("graph id", resp2.GraphId)

	fmt.Println("stream tasks:")
	// TODO: get initial to avoid race condition
	resp3, err3 := client.StreamTasks(ctx, &proto.StreamTasksRequest{
		GraphId:        resp2.GraphId,
		IncludeInitial: true,
	})
	if err3 != nil {
		log.Fatal(err3)
	}

	for {
		evt, err := resp3.Recv()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("evt:", evt)
	}
}
