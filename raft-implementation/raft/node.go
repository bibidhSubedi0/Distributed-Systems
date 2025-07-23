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
	"google.golang.org/grpc/credentials/insecure"
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
	go n.listen(s)
	time.Sleep(time.Millisecond)
	go n.request(s)
	wg.Wait()
}

func (n *Node) Stop() {
	close(n.stopCh)
}

func (n *Node) listen(s *grpc.Server) {
	lis, err := net.Listen("tcp", n.Address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("%s Listner Running\n", n.Address)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

func (n *Node) request(s *grpc.Server) {
	// Make a vote request
	var toconnect = "localhost:5002"
	conn, err := grpc.Dial(toconnect, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	log.Printf("connected to: %s\n", toconnect)
	defer conn.Close()

	client := pb.NewTestServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	rsp, err := client.TestThis(ctx, &pb.TestRequest{Input: "Some input from client", Id: n.ID})
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	} else {
		log.Printf("Response : %s", rsp)
	}

}

func (s *RaftNode) TestThis(ctx context.Context, req *pb.TestRequest) (*pb.TestResponse, error) {
	log.Print("Test requested by : ")
	name := req.GetInput()
	resp := fmt.Sprintf("Hello k xa %s", name)
	return &pb.TestResponse{Resp: resp}, nil
}

func (s *RaftNode) RequestVote(ctx context.Context, req *pb.RequestVoteRequest) (*pb.RequestVoteResponse, error) {
	log.Printf("Vote requested by : %s", s.ID)
	var vote = false

	return &pb.RequestVoteResponse{VoteGiven: vote}, nil
}
