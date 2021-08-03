package net

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "github.com/4gatepylon/GoPoker/net/proto"
)

const dialTimeout = 3
const pingTimeout = 3

func RunClient() {
	conn, err := grpc.Dial(serverAddr(), grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(dialTimeout*time.Second))
	if err != nil {
		log.Fatalf("Failed to dial: %v\n", err)
	}
	defer conn.Close()

	client := pb.NewPingerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout*time.Second)
	defer cancel()

	resp, err := client.Ping(ctx, &pb.PingRequest{})
	if err != nil {
		log.Fatalf("Failed to ping: %v\n", err)
	}
	log.Printf("Response: %s\n", resp.GetMessage())
}
