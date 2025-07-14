package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	node "github.com/bibidhSubedi0/raft/raft"
)

type ConfigNode struct {
	ID      int    `json:"id"`
	Address string `json:"address"`
}

func main() {

	// Read the config file, aile lai just a json
	data, err := os.ReadFile("cluster.json")
	if err != nil {
		log.Fatalf("Failed to read config file : %v", err)
	}

	fmt.Print(data)

	// Unmarsall will put the data in proper fomat of the config
	var config []ConfigNode
	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatalf("Failed to parse config : %v", err)
	}

	// Create the nodes from info retrived from the config file
	nodes := make([]*node.Node, len(config))

	for i, cfg := range config {
		nodes[i] = &node.Node{
			ID:        cfg.ID,
			Address:   cfg.Address,
			Role:      node.Follower,
			Term:      0,
			Log:       []string{},
			VotedFor:  nil,
			Neighbors: []*node.Node{},
		}
	}

	// Assign the required neighbours
	for _, n := range nodes {
		for _, other := range nodes {
			if other.ID != n.ID {
				n.Neighbors = append(n.Neighbors, other)
			}
		}
	}

	// Print the detial of the neighbours
	for _, n := range nodes {
		node.PrintDetails(*n)
		fmt.Println("-----")
	}

}
