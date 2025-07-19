package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	pb "github.com/bibidhSubedi0/raft/proto"
	node "github.com/bibidhSubedi0/raft/raft"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ConfigNode struct {
	ID      string `json:"id"`
	Address string `json:"address"`
}

func main() {

	var inp bool

	_, err := fmt.Scan(&inp)
	if err != nil {
		return
	}

	if inp {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		s := grpc.NewServer()
		pb.RegisterTestServiceServer(s, &server{})

		log.Println("Server is running on :50051")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	} else {
		conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}

		log.Printf("\nconnected")
		defer conn.Close()

		client := pb.NewTestServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		rsp, err := client.TestThis(ctx, &pb.TestRequest{Input: "Some input from client"})
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		} else {
			log.Printf("Response : %s", rsp)
		}

	}

	/*
		nodes := initializeNodes()
		var wg sync.WaitGroup
		wg.Add(len(nodes))
		for _, n := range nodes {
			go func(n *node.Node) {
				defer wg.Done()
				activate(n)
			}(n)
		}
		wg.Wait() // Wait for all goroutines to finish
	*/
}

func initializeNodes() []*node.Node {
	// Read the config file
	data, err := os.ReadFile("cluster.json")
	if err != nil {
		log.Fatalf("\nFailed to read config file: %v", err)
	}

	// Parse config
	var config []ConfigNode
	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatalf("\nFailed to parse config: %v", err)
	}

	// Create nodes
	nodes := make([]*node.Node, len(config))

	// create all node instances
	for i, cfg := range config {
		nodes[i] = node.NewNode(cfg.ID, cfg.Address)
	}

	// set up neighbors for each node
	for _, n := range nodes {
		for _, other := range nodes {
			if other.ID != n.ID {
				node.AddNeighbor(n, other)
			}
		}
	}

	return nodes
}

func activate(n *node.Node) {
	log.Printf("Starting node %s at address %s\n", n.ID, n.Address)

	if err := n.Start(); err != nil {
		log.Printf("Error starting node %s: %v\n", n.ID, err)
		return
	}

}

type server struct {
	pb.UnimplementedTestServiceServer
}

func (s *server) TestThis(ctx context.Context, req *pb.TestRequest) (*pb.TestResponse, error) {
	log.Print("Test requested by : ")
	name := req.GetInput()
	resp := fmt.Sprintf("Hello k xa %s", name)
	return &pb.TestResponse{Resp: resp}, nil
}
