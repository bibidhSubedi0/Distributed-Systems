package node

import (
	"sync"
	"time"
)

type Role int

const (
	Follower Role = iota
	Candidate
	Leader
)

func (r Role) String() string {
	switch r {
	case Follower:
		return "Follower"
	case Candidate:
		return "Candidate"
	case Leader:
		return "Leader"
	default:
		return "Unknown"
	}
}

type Node struct {
	ID        string
	Address   string
	Role      Role
	Term      float32
	VotedFor  *Node
	Neighbors []*Node

	// Timing
	electionTimeout time.Duration
	heartbeatTicker *time.Ticker
	electionTimer   *time.Timer

	// Logs
	Log         []string
	commitIndex int
	lastApplied int

	// Leader-specific State
	nextIndex  map[string]int
	matchIndex map[string]int

	Mu          sync.Mutex    // Protects access to shared state
	stopCh      chan struct{} // Channel for stopping the node
	heartbeatCh chan bool     // Channel for heartbeat events
	electionCh  chan bool     // Channel for election events
}

// NewNode creates and initializes a new Node
func NewNode(id, address string) *Node {
	n := &Node{
		ID:          id,
		Address:     address,
		Role:        Follower,
		Term:        0,
		Log:         make([]string, 0),
		commitIndex: -1,
		lastApplied: -1,
		nextIndex:   make(map[string]int),
		matchIndex:  make(map[string]int),
		stopCh:      make(chan struct{}),
		heartbeatCh: make(chan bool),
		electionCh:  make(chan bool),
	}
	return n
}

func AddNeighbor(n *Node, neighbor *Node) {
	n.Mu.Lock()
	defer n.Mu.Unlock()
	n.Neighbors = append(n.Neighbors, neighbor)
}

func (n *Node) Start() error {
	go n.run()
	return nil
}

func (n *Node) Stop() {
	close(n.stopCh)
}

func (n *Node) run() {
	for {

	}
}
