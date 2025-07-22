package node

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	pb "github.com/bibidhSubedi0/raft/proto"
	"google.golang.org/grpc"
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

type RaftNode struct {
	pb.UnimplementedTestServiceServer
	pb.UnimplementedRequestVoteServiceServer
	*Node
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

func (n *Node) Start(wg *sync.WaitGroup, s *grpc.Server) {
	defer wg.Done()
	fmt.Printf("%s Listner Running\n", n.ID)
	n.listen(s)
}

func (n *Node) Stop() {
	close(n.stopCh)
}

func (n *Node) listen(s *grpc.Server) {
	fmt.Println("Listening")

	lis, err := net.Listen("tcp", n.Address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Server is running on :%s", n.Address)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (n *Node) RequestVote() {

}

func (s *RaftNode) TestThis(ctx context.Context, req *pb.TestRequest) (*pb.TestResponse, error) {
	log.Print("Test requested by : ")
	name := req.GetInput()
	resp := fmt.Sprintf("Hello k xa %s", name)
	return &pb.TestResponse{Resp: resp}, nil
}
