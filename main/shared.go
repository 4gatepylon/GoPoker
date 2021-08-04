package main

import (
	"github.com/4gatepylon/GoPoker/poker"
)

// Requests are a shared type between client and server.
// clients have NetClients which send requests over the network
// in the format specified by Request here. Servers have NetServers
// that similarly recieve Request requests and feed them to the
// parent server to handle.

// To allow request/response batching (lowers network utilization)
// bitshifted integers are used as the request/response type.
// Request/response data is formatted to be able to house any request/response
// but only the relevant portions will be read. Response status will mirror
// request types to allow multiplexxed responses. A one at the same bit will
// represent OK and a zero will present FAIL. A more descriptive error message
// can be left in the data.

// You can imagine UIs as having a button per move.

const (
	NET_RQTYPE_PLAYER_CHECK uint64 = 1 << iota // = poker.MTYPE_CHECK
	NET_RQTYPE_FOLD                            // = poker.MTYPE_FOLD
	NET_RQTYPE_CALL                            // = poker.MTYPE_CALL
	NET_RQTYPE_CALL_ANY                        // = poker.MTYPE_CALL_ANY
	NET_RQTYPE_BET                             // = poker.MTYPE_BET
	NET_RQTYPE_SITOUT_NEXT_ROUND               // = poker.MTYPE_SITOUT_NEXT_ROUND
	NET_RQTYPE_REQ_MOD
	NET_RQTYPE_MESSAGE
	NET_RQTYPE_JOIN
	NET_RQTYPE_LEAVE
	NET_RQTYPE_SHOW_CARDS
	NET_RQTYPE_SHOW_LEFT_CARD
	NET_RQTYPE_SHOW_RIGHT_CARD
)

const (
	NET_RPTYPE_OK uint64 = 1 << iota
	NET_RPTYPE_FAIL
	NET_RPTYPE_UNAUTH_RQ
	NET_RPTYPE_INVALID_RQ
	NET_RPTYPE_ERR
)

const (
	UI_RQTYPE_PLAYER_CHECK uint64 = 1 << iota // = NET_* = poker.MTYPE_CHECK
	UI_RQTYPE_FOLD                            // = NET_* = poker.MTYPE_FOLD
	UI_RQTYPE_CALL                            // = NET_* = poker.MTYPE_CALL
	UI_RQTYPE_CALL_ANY                        // = NET_* = poker.MTYPE_CALL_ANY
	UI_RQTYPE_BET                             // = NET_* = poker.MTYPE_BET
	UI_RQTYPE_SITOUT_NEXT_ROUND               // = NET_* = poker.MTYPE_SITOUT_NEXT_ROUND
	UI_RQTYPE_REQ_MOD                         // = NET_*
	UI_RQTYPE_MESSAGE                         // = NET_*
	UI_RQTYPE_JOIN                            // = NET_*
	UI_RQTYPE_LEAVE                           // = NET_*
	UI_RQTYPE_SHOW_CARDS                      // = NET_*
	UI_RQTYPE_SHOW_LEFT_CARD                  // = NET_*
	UI_RQTYPE_SHOW_RIGHT_CARD                 // = NET_*
	UI_RQTYPE_LOBBY_LIST
	UI_RQTYPE_CLIENT_EXIT
)

const (
	// Lobby Updates
	UI_RPTYPE_LOBBY_LIST_ADD uint64 = (iota + 1)
	UI_RPTYPE_LOBBY_LIST_DELETE

	// Game Updates
	UI_RPTYPE_GAME_PLAYER_NEW                  // A new player (name, id, chips) joins
	UI_RPTYPE_GAME_PLAYER_LEAVE                // An old player (id) leaves

	UI_RPTYPE_GAME_POT_NEW                     // A new pot (chips) is created 
	UI_RPTYPE_GAME_POT_UPDATE                  // A pot (id, value) is updated in value (value = 0 is delete)

	UI_RPTYPE_GAME_MID_ADD                     // Add a card (card) to the middle
	UI_RPTYPE_GAME_MID_CLEAR                   // Clear the middle ()

	UI_RPTYPE_GAME_PLAYER_CHIP_UPDATE          // Update a player's (id, chips) chip count
	UI_RPTYPE_GAME_PLAYER_NAME_UPDATE          // Update a player's (id, name) name
	UI_RPTYPE_GAME_PLAYER_MOD_UPDATE           // Update a player's (id, perms) permissions
	UI_RPTYPE_GAME_PLAYER_STATUS_UPDATE        // Update a player's status (id, status)
	UI_RPTYPE_GAME_PLAYER_POT_UPDATE           // Update a player's pot value (how much they have bet) (id, value)
	UI_RPTYPE_GAME_PLAYER_CARD_LEFT            // Update a player's left card (id, card)
	UI_RPTYPE_GAME_PLAYER_CARD_RIGHT           // Update a player's right card (id, card)
	
	UI_RPTYPE_GAME_GMODE_UPDATE                // Update the game mode (game mode)
	UI_RPTYPE_GAME_GSTATUS_UPDATE              // Update the game status (game status)
	UI_RPTYPE_GAME_BROUND_UPDATE               // Update the betting round (betting round)
	UI_RPTYPE_GAME_ORDER_UPDATE                // Update a single player's location in the playing order (id, location)
	UI_RPTYPE_GAME_STAKES_UPDATE               // Update the stakes (new stakes)

    UI_RPTYPE_GAME_MESSAGE                     // Give the player a string message (message)
	UI_RPTYPE_SHOW_LEFT                        // Show a player's (id) left card
	UI_RPTYPE_SHOW_RIGHT                       // Show a player's (id) right card
	// DEBUG: 22 messages above

	// Control Messages
	UI_RPTYPE_UPDATE_ENACT_OK uint64 = 1 << 63 // Pass true/one normally; pass false/zero to await next true/one to enact change
)

type UIRequest struct {
	TypeInt uint64
	NameStr *string
	MsgStr  *string
}

type UIResponse struct {
	Type   uint64
	ValInt uint64
	IdInt  uint64
	Str    *string
}

type NetRequest struct {
	Type      uint64
	RequestId uint64
	NameStr   *string
	MsgStr    *string
}

// The zero request id is reserved for server-initiated messages
// (i.e. updates to the game state, etc...).
type NetResponse struct {
	Type      uint64
	RequestId uint64

	// Data includes any and all game state information necessary to update it.
	Middle  [5]*poker.CardLike
	Players []*poker.PlayerInfo
	Pots    []uint64
}

type NetServerLike interface {
	// External Control
	Serve(*string) (chan *NetRequest, chan *NetResponse, error)                        // (local addr) => (reqs, resps, error)
	Close() error                                                                      // () => (error)

	// Internallly it should be listening to requests from the network and then
	// parsing them and converting them into *Request objects to send on the channel.
	// At the same time it should recieve response objects on the other channel to
	// format into network types and send on the wire.
}

type ServerLike interface {
	Init(chan *NetRequest, chan *NetResponse) error                                    // (reqs, resps) => (error)
	Start() error                                                                      // () => (error)
	Stop() error                                                                       // () => (error)
	Teardown() error                                                                   // () => (error)

	// Internally Init should create any necessary metadata and files (etc)
	// while start should spin up the server, stop should stop the server (without
	// deleting any information on disk if it exists) and teardown should tear the server
	// down cleaning up after itself. While it is running, it should listen for requests on
	// the first channel and issue responses on the second channel. You can think of a
	// Server as containing an imaginary game server (the main thread or process or machine)
	// that manages various games, and a network server that listens for requests. These two
	// communicate via the two channels.
}

type NetClientLike interface {
	Connect(*string) (chan *NetRequest, chan *NetResponse, error)                      // (remote addr) => (reqs, resps, error)
	Close() error                                                                      // () => (error)

	// Internally, a net client is like the net server in that it encodes/decodes *Request
	// formatted request into network types to send on the wire, but instead of serving,
	// it connects to a server.
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

type UILike interface {
	Init() (chan *UIRequest, chan *UIResponse, error)                                  // () => (ui reqs, ui resps, error)
	Start() error                                                                      // () => error
	Stop() error                                                                       // () => error

	// A UI should be a CLI or a GUI. For UX it should probably have additional settings for
	// color schemes, keybindings, and other such things.
}