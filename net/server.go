package net

import (
	"context"
	"log"
	"net"
	"fmt"
	"io"

	"google.golang.org/grpc"
	"github.com/4gatepylon/GoPoker/protocol"
	"github.com/4gatepylon/GoPoker/poker"

	pb "github.com/4gatepylon/GoPoker/net/proto"
)

// Server implement both protocol.ServerLike and protocol.NetServerLike
type Server struct {
	pb.UnimplementedGameServerServer
	requests          chan *protocol.NetRequest
	responses         chan *protocol.NetResponse
	Addr              string
	running           bool
	games             map[uint64]poker.GameLike
	grpcServer        *grpc.Server
}

func (server *Server) JoinGame(ctx context.Context, joinReq *pb.JoinGameRequest) (*pb.JoinGameResponse, error) {
	// TODO recieve a request, send some back and forth on the channel, then send a response
	return nil, nil // this is basically a rest request, we will ignore the channels here
}
func (server *Server) CreateGame(ctx context.Context, createReq *pb.CreateGameRequest) (*pb.CreateGameResponse, error) {
	// TODO same as above
	return nil, nil // this is basically a rest request, we will ignore the channels here
}
func (server *Server) LeaveGame(ctx context.Context, leaveReq *pb.LeaveGameRequest) (*pb.LeaveGameResponse, error) {
	// TODO ibid
	return nil, nil // this is basically a rest request, we will ignore the channels here
}
func (server *Server) GameStream(stream pb.GameServer_GameStreamServer) error {
	// Handle incoming requests from the network and issue them on the channel
	go func() {
		for {
			_, err := stream.Recv()
			if err == io.EOF {
				log.Printf("Client closed the connection\n")
				return
			}
			if err != nil {
				return // TODO
			}
			// do translation TODO
			server.requests <- &protocol.NetRequest{
				Type: protocol.NET_RQTYPE_CHECK,
			}
		}
	}()

	// Handle outgoing responses from the channel, form that them, and send on the net
	go func() {
		for response := range server.responses {
			if response != nil {
				log.Printf("Yay!\n")
			}

			if response.Type == protocol.NET_RPTYPE_ERROR {
				log.Printf("Just kidding, error!\n")
			}
		}
	}()

	// Handle incoming requests from the channel and handle them for the games, then issue
	// the appropriate responses
	for request := range server.requests {
		if request.Type | protocol.NET_RQTYPE_CHECK > 0 {

		}
		if request.Type | protocol.NET_RQTYPE_FOLD > 0 {

		}
		if request.Type | protocol.NET_RQTYPE_CALL > 0 {

		}
		if request.Type | protocol.NET_RQTYPE_CALL_ANY > 0 {

		}
		if request.Type | protocol.NET_RQTYPE_BET > 0 {

		}
		if request.Type | protocol.NET_RQTYPE_SITOUT_NEXT_ROUND > 0 {

		}
		if request.Type | protocol.NET_RQTYPE_REQ_MOD > 0 {

		}
		if request.Type | protocol.NET_RQTYPE_MESSAGE > 0 {

		}
		if request.Type | protocol.NET_RQTYPE_JOIN > 0 {

		}
		if request.Type | protocol.NET_RQTYPE_CREATE > 0 {

		}
		if request.Type | protocol.NET_RQTYPE_LEAVE > 0 {

		}
		if request.Type | protocol.NET_RQTYPE_SHOW_CARDS > 0 {

		}
		if request.Type | protocol.NET_RQTYPE_SHOW_LEFT_CARD > 0 {

		}
		if request.Type | protocol.NET_RQTYPE_SHOW_RIGHT_CARD > 0 {

		}
	}
	return nil
}

func (server *Server) Start() error {
	// TODO send a "starting" message to connections awaiting restart
	server.running = true;
	return nil
}

func (server *Server) Stop() error {
	// TODO send a "stopping" message
	server.running = false;
	return nil
}

func (server *Server) Teardown() error {
	for stream_code, game := range server.games {
		// TODO send a "tearing down" message
		close(server.requests)
		close(server.responses)
		err := game.Teardown()
		if err != nil {
			return fmt.Errorf("Failed to shut down game with stream code %d and err `%v`\n", stream_code, err)
		}
	}

	return nil
}

// Not Used
func (server *Server) Init(rq chan *protocol.NetRequest, rp chan *protocol.NetResponse) error {
	return fmt.Errorf("Not used in this server implementation\n")
}

func (server *Server) Close() error {
	return fmt.Errorf("Not used in this server implementation\n")
}

func (server *Server) Serve(*string) (chan *protocol.NetRequest, chan *protocol.NetResponse, error) {
	return nil, nil, fmt.Errorf("Serve is not used in this implementation\n")
}

// Proper flow is:
// (1) create struct and any files/procs necessary
// (2) Start
// (3) Stop 
// (4) Teardown
// NOTE: Close(), Init(), and Serve() are all ignored in this implementation. They are meant for implementations
// that are split into different structs and need communication and initialization for both types.
func NewServer(addr string) (*Server, error) {
	s := &Server{
		requests:          make(chan *protocol.NetRequest, 8),
		responses:         make(chan *protocol.NetResponse, 8),
		Addr:              addr,
		running:           true,
		games:             make(map[uint64]poker.GameLike, 0),
	}

	lis, err := net.Listen("tcp", serverAddr())
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterGameServerServer(grpcServer, s)

	log.Printf("server listening at %v\n", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		return nil, fmt.Errorf("failed to serve: `%v`\n", err)
	}

	return s, nil
}