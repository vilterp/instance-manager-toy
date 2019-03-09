package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"

	"github.com/cockroachlabs/instance_manager/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var rootCmd = cobra.Command{
	Use: "client [command]",
}

var nodesCmd = &cobra.Command{
	Use: "nodes",
}

var nodesLsCommand = &cobra.Command{
	Use:  "ls",
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		resp, err := c.ListNodes(context.Background(), &proto.ListNodesRequest{})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for _, node := range resp.Instances {
			fmt.Println(node)
		}
	},
}

var nodesStreamCommand = &cobra.Command{
	Use:  "stream",
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		resp, err := c.StreamNodes(context.Background(), &proto.StreamNodesRequest{
			IncludeInitial: true,
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for {
			evt, err := resp.Recv()
			if err != nil {
				fmt.Println("err:", err)
				return
			}
			fmt.Println(evt)
		}
	},
}

var graphsCmd = &cobra.Command{
	Use: "graphs",
}

var graphsLsCmd = &cobra.Command{
	Use:  "ls",
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		resp, err := c.ListTaskGraphs(context.Background(), &proto.ListTaskGraphsRequest{})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("ID\tState\tCreated At\tStarted At\tFinished At")
		for _, tg := range resp.TaskGraphs {
			fmt.Printf(
				"%v\t%v\t%v\t%v\t%v\n",
				tg.Id, tg.State, formatTimestamp(tg.CreatedAt),
				formatTimestamp(tg.StartedAt), formatTimestamp(tg.FinishedAt),
			)
		}
	},
}

var graphsGetCmd = &cobra.Command{
	Use:  "get",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		graphID := args[0]
		resp, err := c.GetTaskGraph(context.Background(), &proto.GetTaskGraphRequest{Id: graphID})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		graph := resp.Graph
		spec := graph.Spec
		tasks := graph.Tasks
		graph.Spec = nil
		graph.Tasks = nil

		fmt.Println(graph)
		fmt.Println()
		fmt.Println("SPEC:")
		spec.Print()
		fmt.Println()
		fmt.Println("TASKS:")
		sort.Slice(tasks, func(i, j int) bool {
			return tasks[i].Id < tasks[j].Id
		})
		fmt.Println("ID\tAction\tStartedAt\tFinishedAt\tErr")
		for _, t := range tasks {
			fmt.Printf(
				"%v\t%v\t%v\t%v\t%v\n",
				t.Id, t.Action, formatTimestamp(t.StartedAt),
				formatTimestamp(t.FinishedAt), t.Error,
			)
		}
	},
}

var tasksStreamCmd = &cobra.Command{
	Use:  "stream-tasks",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		graphID := args[0]
		resp, err := c.StreamTasks(context.Background(), &proto.StreamTasksRequest{
			IncludeInitial: true,
			GraphId:        graphID,
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for {
			evt, err := resp.Recv()
			if err != nil {
				log.Println(err)
				return
			}
			fmt.Println(evt)
		}
	},
}

var updateCommand = &cobra.Command{
	Use:  "update [num] [version]",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		num, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		vers, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		update(&proto.GroupSpec{
			NumInstances: int64(num),
			Version:      int64(vers),
		})
	},
}

//var graphsStreamCmd = &cobra.Command{
//	Use: "stream",
//	Run: func(cmd *cobra.Command, args []string) {
//		c := getClient()
//		resp, err := c.StreamTaskGraphs(context.Background(), &proto.StreamTaskGraphsRequest{})
//	},
//}

func main() {
	rootCmd.AddCommand(updateCommand)
	rootCmd.AddCommand(nodesCmd)
	nodesCmd.AddCommand(nodesLsCommand)
	nodesCmd.AddCommand(nodesStreamCommand)
	rootCmd.AddCommand(graphsCmd)
	graphsCmd.AddCommand(graphsLsCmd)
	graphsCmd.AddCommand(graphsGetCmd)
	graphsCmd.AddCommand(tasksStreamCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("fatal error:", err)
		os.Exit(1)
	}

}

func getClient() proto.GroupManagerClient {
	addr := "0.0.0.0:8888"

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("failed to dial:", err)
	}
	client := proto.NewGroupManagerClient(conn)
	return client
}

func update(newSpec *proto.GroupSpec) {
	client := getClient()

	ctx := context.Background()
	resp2, err2 := client.UpdateSpec(ctx, &proto.UpdateSpecRequest{
		Spec: newSpec,
	})
	if err2 != nil {
		log.Fatal(err2)
	}
	fmt.Println("TASK GRAPH SPEC:")
	resp2.Graph.Spec.Print()
	fmt.Println()

	streamTasks(client, ctx, resp2.Graph.Id)
}

func streamTasks(client proto.GroupManagerClient, ctx context.Context, graphID string) {
	fmt.Println("TASKS STREAM:")
	fmt.Println()
	resp, err := client.StreamTasks(ctx, &proto.StreamTasksRequest{
		GraphId:        graphID,
		IncludeInitial: true,
	})
	if err != nil {
		log.Println(err)
		return
	}

	for {
		evt, err := resp.Recv()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("task evt:", evt)
	}
}

func formatTimestamp(ts *timestamp.Timestamp) string {
	return ptypes.TimestampString(ts)
}
