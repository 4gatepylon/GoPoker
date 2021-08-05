package protocol

import (
	"github.com/4gatepylon/GoPoker/poker"
)

// To allow request/response batching (lowers network utilization)
// bitshifted integers are used as the request/response type.
// Request/response data is formatted to be able to house any request/response
// but only the relevant portions will be read. Response status will mirror
// request types to allow multiplexxed responses. A one at the same bit will
// represent OK and a zero will present FAIL. A more descriptive error message
// can be left in the data.

// You can imagine UIs as having a button per move.

const (
	NET_RQTYPE_CHECK  uint64 = 1 << iota            // = poker.MTYPE_CHECK
	NET_RQTYPE_FOLD                                 // = poker.MTYPE_FOLD
	NET_RQTYPE_CALL                                 // = poker.MTYPE_CALL
	NET_RQTYPE_CALL_ANY                             // = poker.MTYPE_CALL_ANY
	NET_RQTYPE_BET                                  // = poker.MTYPE_BET
	NET_RQTYPE_SITOUT_NEXT_ROUND                    // = poker.MTYPE_SITOUT_NEXT_ROUND
	NET_RQTYPE_REQ_MOD
	NET_RQTYPE_MESSAGE
	NET_RQTYPE_JOIN
	NET_RQTYPE_CREATE
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
	NET_RPTYPE_ERROR
)

const (
	UI_RQTYPE_CHECK uint64 = 1 << iota             // = NET_* = poker.MTYPE_CHECK
	UI_RQTYPE_FOLD                                 // = NET_* = poker.MTYPE_FOLD
	UI_RQTYPE_CALL                                 // = NET_* = poker.MTYPE_CALL
	UI_RQTYPE_CALL_ANY                             // = NET_* = poker.MTYPE_CALL_ANY
	UI_RQTYPE_BET                                  // = NET_* = poker.MTYPE_BET
	UI_RQTYPE_SITOUT_NEXT_ROUND                    // = NET_* = poker.MTYPE_SITOUT_NEXT_ROUND
	UI_RQTYPE_REQ_MOD                              // = NET_*
	UI_RQTYPE_MESSAGE                              // = NET_*
	UI_RQTYPE_JOIN                                 // = NET_*
	UI_RQTYPE_LEAVE                                // = NET_*
	UI_RQTYPE_SHOW_CARDS                           // = NET_*
	UI_RQTYPE_SHOW_LEFT_CARD                       // = NET_*
	UI_RQTYPE_SHOW_RIGHT_CARD                      // = NET_*
	UI_RQTYPE_LOBBY_LIST
	UI_RQTYPE_CLIENT_EXIT
)

const (
	// Lobby Updates
	UI_RPTYPE_LOBBY_LIST_ADD uint64 = (iota + 1)
	UI_RPTYPE_LOBBY_LIST_DELETE

	// Game Updates
	UI_RPTYPE_GAME_PLAYER_NEW   // A new player (name, id, chips) joins
	UI_RPTYPE_GAME_PLAYER_LEAVE // An old player (id) leaves

	UI_RPTYPE_GAME_POT_NEW    // A new pot (chips) is created
	UI_RPTYPE_GAME_POT_UPDATE // A pot (id, value) is updated in value (value = 0 is delete)

	UI_RPTYPE_GAME_MID_ADD   // Add a card (card) to the middle
	UI_RPTYPE_GAME_MID_CLEAR // Clear the middle ()

	UI_RPTYPE_GAME_PLAYER_CHIP_UPDATE   // Update a player's (id, chips) chip count
	UI_RPTYPE_GAME_PLAYER_NAME_UPDATE   // Update a player's (id, name) name
	UI_RPTYPE_GAME_PLAYER_MOD_UPDATE    // Update a player's (id, perms) permissions
	UI_RPTYPE_GAME_PLAYER_STATUS_UPDATE // Update a player's status (id, status)
	UI_RPTYPE_GAME_PLAYER_POT_UPDATE    // Update a player's pot value (how much they have bet) (id, value)
	UI_RPTYPE_GAME_PLAYER_CARD_LEFT     // Update a player's left card (id, card)
	UI_RPTYPE_GAME_PLAYER_CARD_RIGHT    // Update a player's right card (id, card)

	UI_RPTYPE_GAME_GMODE_UPDATE   // Update the game mode (game mode)
	UI_RPTYPE_GAME_GSTATUS_UPDATE // Update the game status (game status)
	UI_RPTYPE_GAME_BROUND_UPDATE  // Update the betting round (betting round)
	UI_RPTYPE_GAME_ORDER_UPDATE   // Update a single player's location in the playing order (id, location)
	UI_RPTYPE_GAME_STAKES_UPDATE  // Update the stakes (new stakes)

	UI_RPTYPE_GAME_MESSAGE // Give the player a string message (message)
	UI_RPTYPE_SHOW_LEFT    // Show a player's (id) left card
	UI_RPTYPE_SHOW_RIGHT   // Show a player's (id) right card
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

// Translators are effectively a lightweight interface used by internals of servers and clients to translate
// UIRequest/Response => NetRequest/Response => Protocol buffer defined messages sendable by GRPC. We have an
// interface in case we wish to swap out GRPC with REST or some other framework.

// NetRequestTranslators turn net requests into proto. NetResponseTranslators turn proto into NetResponse.
// UIRequestTranslators turn UIRequests into NetRequests. UIResponseTranslators turn NetResponses into UIResponses.
// It is possible in any of these cases to produce zero requests/responses from one or more messages (i.e. drop).
// It is also possible to accumulate requests/responses and then create joint messages to send on the wire (or channel).

// Internally, consume should be updating some form of state. Produce should be flushing it. That is the same
// for all of these translators.

type NetRequestTranslator interface {
	Consume(*NetRequest) (uint64, error)    // (incoming net req) => (how many consumed, error)
	Produce([]*interface{}) (uint64, error) // (outgoing proto) => (how many produced, error)
}

type UIRequestTranslator interface {
	Consume(*UIRequest) (uint64, error)    // (incoming ui req) => (how many produced, error)
	Produce(*[]NetRequest) (uint64, error) // (outgoing net req) => (how many produced, error)
}

type NetResponseTranslator interface {
	Consume(*interface{}) (uint64, error)   // (incoming proto) => (how many produced, error)
	Produce(*[]NetResponse) (uint64, error) // (outgoing net resp) => (how many produced, error)
}

type UIResponseTranslator interface {
	Consume(*NetResponse) (uint64, error)  // (incoming net resp) => (how many produced, error)
	Produce([]*UIResponse) (uint64, error) // (outgoing ui resp) => (how many produced, error)
}
