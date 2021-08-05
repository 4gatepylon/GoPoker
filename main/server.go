package main

type NetServerLike interface {
	// External Control
	Serve(*string) (chan *NetRequest, chan *NetResponse, error) // (local addr) => (reqs, resps, error)
	Close() error                                               // () => (error)

	// Internallly it should be listening to requests from the network and then
	// parsing them and converting them into *Request objects to send on the channel.
	// At the same time it should recieve response objects on the other channel to
	// format into network types and send on the wire.
}

type ServerLike interface {
	Init(chan *NetRequest, chan *NetResponse) error // (reqs, resps) => (error)
	Start() error                                   // () => (error)
	Stop() error                                    // () => (error)
	Teardown() error                                // () => (error)

	// Internally Init should create any necessary metadata and files (etc)
	// while start should spin up the server, stop should stop the server (without
	// deleting any information on disk if it exists) and teardown should tear the server
	// down cleaning up after itself. While it is running, it should listen for requests on
	// the first channel and issue responses on the second channel. You can think of a
	// Server as containing an imaginary game server (the main thread or process or machine)
	// that manages various games, and a network server that listens for requests. These two
	// communicate via the two channels.
}
