package node

import "fmt"

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
	ID        int
	Address   string
	Role      Role
	Term      float32
	Log       []string
	VotedFor  *Node
	Neighbors []*Node
}

func PrintDetails(node Node) {
	fmt.Printf("ID: %d\n", node.ID)
	fmt.Printf("Address: %s\n", node.Address)
	fmt.Printf("Role: %s\n", node.Role.String())
	fmt.Printf("Term: %.2f\n", node.Term)
	fmt.Printf("Log entries: %v\n", node.Log)

	if node.VotedFor != nil {
		fmt.Printf("VotedFor Node ID: %d\n", node.VotedFor.ID)
	} else {
		fmt.Println("VotedFor Node: nil")
	}

	fmt.Printf("Neighbors IDs: ")
	if len(node.Neighbors) == 0 {
		fmt.Println("none")
	} else {
		for _, neighbor := range node.Neighbors {
			fmt.Printf("%d ", neighbor.ID)
		}
		fmt.Println()
	}
}

func main() {
	node1 := &Node{
		ID:        1,
		Address:   "x",
		Role:      Follower,
		Term:      0,
		Log:       []string{},
		VotedFor:  nil,
		Neighbors: []*Node{},
	}
	PrintDetails(*node1)
}
