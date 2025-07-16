package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pb "test/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedDateServiceServer
}

func (s *server) GetCurrentDate(ctx context.Context, req *pb.DateRequest) (*pb.DateResponse, error) {
	currentTime := time.Now().UTC().Format("2006-01-02 15:04:05")
	return &pb.DateResponse{
		CurrentDate: currentTime,
	}, nil
}

func (s *server) GetUserInfo(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	return &pb.UserResponse{
		UserLogin: "bibidhSubedi0",
	}, nil
}

func (s *server) Greet(ctx context.Context, req *pb.GreetRequest) (*pb.GreetResponse, error) {
	name := req.GetName()
	resp := fmt.Sprintf("Hello k xa %s", name)
	return &pb.GreetResponse{Text: resp}, nil

}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterDateServiceServer(s, &server{})

	log.Println("Server is running on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
