package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	node "github.com/bibidhSubedi0/raft/raft"
)

type ConfigNode struct {
	ID      string `json:"id"`
	Address string `json:"address"`
}

func main() {

	nodes := initializeNodes()

	var wg sync.WaitGroup

	activate(nodes[0])

	for _, n := range nodes {
		wg.Add(1)
		go func(n *node.Node) {
			defer wg.Done()
			activate(n)
		}(n)
	}
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
				AddNeighbor(n, other)
			}
		}
	}

	return nodes
}

func AddNeighbor(n *node.Node, neighbor *node.Node) {
	n.Mu.Lock()
	defer n.Mu.Unlock()
	n.Neighbors = append(n.Neighbors, neighbor)
}

func activate(n *node.Node) {
	log.Printf("Starting node %s at address %s\n", n.ID, n.Address)

	if err := n.Start(); err != nil {
		log.Printf("Error starting node %s: %v\n", n.ID, err)
		return
	}

}
