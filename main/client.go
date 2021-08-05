package main

type UILike interface {
	Init() (chan *UIRequest, chan *UIResponse, error) // () => (ui reqs, ui resps, error)
	Start() error                                     // () => error
	Stop() error                                      // () => error

	// A UI should be a CLI or a GUI. For UX it should probably have additional settings for
	// color schemes, keybindings, and other such things.
}

type ClientLike interface {
	Init(chan *NetRequest, chan *NetResponse, chan *UIRequest, chan *UIResponse) error // (net reqs, net resps, ui reqs, ui resps) => (error)
	Start() error                                                                      // () => (error)
	Stop() error                                                                       // () => (error)

	// Internally, you can think of a client as having a UIClient and a NetClient. The UIClient listens to
	// for the user to press buttons and do other things to request stuff. Then the UIClient sends that on the
	// UIRequest channel. The client parses it and formats it into a requst, forwarding it to the NetClient.
	// The net client then requests from the server. When the server responsds, it's decoded, parsed, re-encoded,
	// and sent to the UI. Not all thing the UI sends may be sent to the server, and not all things the server
	// sends back may be sent to the UI.
}

type NetClientLike interface {
	Connect(*string) (chan *NetRequest, chan *NetResponse, error) // (remote addr) => (reqs, resps, error)
	Close() error                                                 // () => (error)

	// Internally, a net client is like the net server in that it encodes/decodes *Request
	// formatted request into network types to send on the wire, but instead of serving,
	// it connects to a server.
}
