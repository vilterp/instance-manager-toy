package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cockroachlabs/instance_manager/pure_manager/proto"
	"google.golang.org/grpc"
)

func main() {
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
	fmt.Println("\t", resp)

	newSpec := &proto.GroupSpec{
		NumInstances: 3,
		Version:      1,
	}
	fmt.Println("update to", newSpec)
	resp2, err2 := client.UpdateSpec(ctx, &proto.UpdateSpecRequest{
		Spec: newSpec,
	})
	if err2 != nil {
		log.Fatal(err2)
	}
	fmt.Println("task graph spec:")
	resp2.Graph.Spec.Print()

	//go streamTasks(client, ctx, resp2.Graph.Id)
	streamNodes(client, ctx) /**/
}

func streamNodes(client proto.GroupManagerClient, ctx context.Context) {
	fmt.Println("stream nodes:")
	resp, err := client.StreamNodes(ctx, &proto.StreamNodesRequest{
		IncludeInitial: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	for {
		evt, err := resp.Recv()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("\tnode evt:", evt)
	}
}

func streamTasks(client proto.GroupManagerClient, ctx context.Context, graphID string) {
	fmt.Println("stream tasks:")
	resp3, err3 := client.StreamTasks(ctx, &proto.StreamTasksRequest{
		GraphId:        graphID,
		IncludeInitial: true,
	})
	if err3 != nil {
		log.Fatal(err3)
	}

	for {
		evt, err := resp3.Recv()
		if err != nil {
			log.Fatalf("%#v", err)
		}
		fmt.Println("\ttask evt:", evt)
	}
}
