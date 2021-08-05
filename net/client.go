package net

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"github.com/4gatepylon/GoPoker/protocol"

	pb "github.com/4gatepylon/GoPoker/net/proto"
)

// Client implements NetClient and Client, but the UI is done seperately
// so that we can more easily add something like SDL in the future. Moreover,
// the language differential is bigger.
type Client struct {
	uiRequests   chan *protocol.UIRequest
	uiResponses  chan *protocol.UIResponse
	netRequests  chan *protocol.NetRequest
	netResponses chan *protocol.NetResponse
	ServerAddr   string
}

func (client *Client) JoinGame(ctx context.Context, in *pb.JoinGameRequest, opts ...grpc.CallOption) (*pb.JoinGameResponse, error) {
	return nil, nil // TODO
}
func (client *Client) CreateGame(ctx context.Context, in *pb.CreateGameRequest, opts ...grpc.CallOption) (*pb.CreateGameResponse, error) {
	return nil, nil // TODO
}
func (client *Client) LeaveGame(ctx context.Context, in *pb.LeaveGameRequest, opts ...grpc.CallOption) (*pb.LeaveGameResponse, error) {
	return nil, nil // TODO
}
func (client *Client) GameStream(ctx context.Context, opts ...grpc.CallOption) (pb.GameServer_GameStreamClient, error) {
	return nil, nil // TODO
}

func (client *Client)  Close() error {
	return nil // TODO
}

// Not Used
func (client *Client) Connect(serverAddr *string) (chan *protocol.NetRequest, chan *protocol.NetResponse, error) {
	return nil, nil, nil
}

func (client *Client) Init(netRq chan *protocol.NetRequest, netRp chan *protocol.NetResponse, uiRq chan *protocol.UIRequest, uiRp chan *protocol.UIResponse) error {
	return nil
}

func (client *Client) Start() error {
	return nil
}

func (client *Client) Stop() error {
	return nil
}

// Proper client flow is:
// (1) Create any necessary structs, files, etc... (and do server connection, etc...)
// (2) Close()
// NOTE: None of the other methods are used because there is no need to do so and this is a simple client.
// I may go ahead and simplify the design going forward TODO.
func NewClient(serverAddr string) (*Client, error) {
	c := &Client{
		uiRequests:   make(chan *protocol.UIRequest, 8),
		uiResponses:  make(chan *protocol.UIResponse, 8),
		netRequests:  make(chan *protocol.NetRequest, 8),
		netResponses: make(chan *protocol.NetResponse, 8),
		ServerAddr:   serverAddr,
	}

	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(3*time.Second))
	if err != nil {
		log.Fatalf("Failed to dial: %v\n", err)
	}
	defer conn.Close()

	return c, nil // TODO
}
