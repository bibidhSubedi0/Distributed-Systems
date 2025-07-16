package main

import (
	"context"
	"log"
	"time"

	pb "test/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewDateServiceClient(conn)

	// Get current date
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	dateResp, err := client.GetCurrentDate(ctx, &pb.DateRequest{})
	if err != nil {
		log.Fatalf("could not get date: %v", err)
	}
	log.Printf("Current Date: %s", dateResp.CurrentDate)

	// Get user info
	userResp, err := client.GetUserInfo(ctx, &pb.UserRequest{})
	if err != nil {
		log.Fatalf("could not get user info: %v", err)
	}
	log.Printf("User Login: %s", userResp.UserLogin)

	// Get greeting
	rsp, err := client.Greet(ctx, &pb.GreetRequest{Name: "aalu paratha"})
	if err != nil {
		log.Fatalf("Cant greet %v", err)
	}
	log.Printf("Greeting : %s", rsp)
}
