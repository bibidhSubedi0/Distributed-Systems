package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	pb "github.com/bibidhSubedi0/raft/proto"
	node "github.com/bibidhSubedi0/raft/raft"
	"google.golang.org/grpc"
)

type ConfigNode struct {
	ID      string `json:"id"`
	Address string `json:"address"`
}

func main() {

	// Initlize the nodes
	nodes := initializeNodes()

	rn := &node.RaftNode{}

	// Register the required services
	s := grpc.NewServer()
	pb.RegisterTestServiceServer(s, rn)
	// pb.RegisterRequestVoteServiceServer(s, rn)

	// Activate the nodes
	var wg sync.WaitGroup
	for _, n := range nodes {
		wg.Add(1)
		activate(n, &wg, s)
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

func activate(n *node.Node, wg *sync.WaitGroup, s *grpc.Server) {
	go n.Start(wg, s)
}
