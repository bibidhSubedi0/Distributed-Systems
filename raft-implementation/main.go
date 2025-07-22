package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
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

	// Initlize the nodes
	nodes := initializeNodes()

	// Activate the nodes
	var wg sync.WaitGroup
	for _, n := range nodes {
		wg.Add(1)
		activate(n, &wg)
	}

	// Now i have all the nodes listening

	wg.Wait() // Will wait until all goroutines call Done()
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

func activate(n *node.Node, wg *sync.WaitGroup) {
	log.Printf("Starting node %s at address %s\n", n.ID, n.Address)
	go n.Start(wg)

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

func listen() {
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
}

func request() {
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
