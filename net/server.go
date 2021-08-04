package net

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/4gatepylon/GoPoker/net/proto"
)

type NetServer struct {
	// NOTE: this is an embedding (i.e. inherits methods and values) and not just
	// a field.
	pb.UnimplementedPingerServer
}

func (s *NetServer) Ping(context.Context, *pb.PingRequest) (*pb.PingReply, error) {
	log.Printf("Recieved ping request\n")
	return &pb.PingReply{
		Message: "OK",
	}, nil
}

func RunServer() {
	lis, err := net.Listen("tcp", serverAddr())
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}
	server := grpc.NewServer()
	pb.RegisterPingerServer(server, &NetServer{})
	log.Printf("Listening at %v\n", lis.Addr())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}
}
