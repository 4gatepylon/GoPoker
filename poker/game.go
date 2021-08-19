package poker

// Every game is broken into an infinite sequence of rounds which is broken up into four betting rounds.
// In any given round, players are playing or not, and in any game they are admins or not.
// Playing players can check, fold, call, bet, call any (plans a future action) or sit out the next
// round (plans a future action). A game can be playing or not (paused) and private or not (public).
// The only currently supported game mode is regular no-limit Texas hold'em with no rake.

// Games are joinable by join codes (passwords) if private, and simply by request if public. Inside each
// game players can become admins or lose their admin status. The creator of a game is the original admin.
// Admins can give out chips and change the game flow (within limits) by pausing (etc).

// If batch moves are supported (i.e. multiple moves at once like "check/call") then in the case of
// potentially nonsensical combinations (like fold/bet) an error should be thrown. Moves should occur
// by presdence: check > fold > call > call any > bet > sit out next round. Presedence is currently encoded
// by numerical value (lowest value is highest presedence). The null move does nothing.

// Batch requests for player status, game status, and game mode should may also be supported. Currently,
// a single status object (being a number) shares multiple orthogonal statuses (i.e. private | playing)
// and because they are independent there is no precedence defined.

// Betting Rounds
const (
	BROUND_PREFLOP uint64 = (iota + 1)
	BROUND_FLOP
	BROUND_TURN
	BROUND_RIVER
)

// Move Types
const (
	MTYPE_CHECK uint64 = 1 << iota
	MTYPE_FOLD
	MTYPE_CALL
	MTYPE_CALL_ANY
	MTYPE_BET
	MTYPE_SITOUT_NEXT_ROUND
)

// Player Status
const (
	PSTATUS_ADMIN uint64 = 1 << iota
	PSTATUS_PLAYING
)

// Player Permissions (right now same as status)
const (
	PPERM_ADMIN uint64 = 1 << iota
	_
)

// Game Status
const (
	GSTATUS_PLAYING uint64 = 1 << iota
	GSTATUS_PRIVATE
)

// Game Mode
const (
	GMODE_CONST_STAKES uint64 = 1 << iota
)

// Defaults for standard games
const (
	DEFAULT_STAKES uint64                 = 1000               // Default big blind is 1000 chips
	DEFAULT_MAX_PLAYERS uint64            = 6                  // By default allow six players maximum
	DEFAULT_MODE uint64                   = GMODE_CONST_STAKES // Standard games have constant stakes.
	DEFAULT_STATUS uint64                 = 0                  // The default status is the null status (not started)
	DEFAULT_STAKES_HAND_MULTIPLIER uint64 = 100                // DEFAULT_STAKES * ..._MULTIPLIER = default starting hand
)

// A GameLike should be able to manipulate CardLikes accordingly. The string method
// will be desired to communiate with players. Format is "<number><suit>" i.e. "10H" for ten of hearts.
type CardLike interface {
	String() string
}

type PlayerInfo struct {
	// Human-Identifiers
	Name  string
	Chips uint64
	Bet   uint64

	// Game Server-Only
	Id    uint64
	Cards [2]CardLike
	Mod   bool
}

// In game likes, control plane functions are used by game servers
// to control the flow of the game at the request of users. Regular control
// functions are usually triggered by specific player requests.
type GameLike interface {
	// Player Control
	KickPlayer(*string, *string) (bool, error)        // (kicker name, kicked name) => (kicked, error)
	ModPlayer(*string, *string, uint64) (bool, error) // (modder name, modded name, mod perms) => (modded, error)

	// Player Control Plane
	AddPlayer(*string, *string) (*string, bool, error) // (prospective player name, join code) => (player name, joined, error)
	Players() []*PlayerInfo                            // () => (an informative list of players in order of play)
	Stakes() uint64                                    // () => (value of big blind in chips)
	Middle() *[5]CardLike                              // () => (array of cards in the middle)
	Pots() []uint64                                    // () => (a slice of monetary values of pots)

	// Game Status
	ChangeGameName(*string, *string) (bool, error) // (name changer, desired name) => (changed name, error)
	Play(*string) (bool, error)                    // (play requester) => (played, error)
	Pause(*string) (bool, error)                   // (pause requester) => (paused, error)
	MakePrivate(*string) (bool, error)             // (make private requester) => (made private, error)
	MakePublic(*string) (bool, error)              // (make public requester) => (made public, error)
	Playing() bool                                 // () => (game is playing)
	Private() bool                                 // () => (game is private)

	// Game Flow
	Move(uint64, uint64, *string) (bool, error)                        // (move, chips: optional, mover) => (moved, error)
	ChangePlayerName(*string, *string, *string) (*string, bool, error) // (namer, player, new name) => (new name, renamed, error)
	GiveChips(*string, *string, uint64) (bool, error)

	// Game Flow Control Plane
	Increment() (bool, error)  // () => (incremented, error)
	Resolve() (*string, error) // () => (winners' informative message, error)
	NewRound() error           // () => (error)
	Renew() error              // () => (error)

	// System Maintenance
	// There should exist a function NewGame(...) or InitGame(...) that
	// creates a new one and initializes any side-effects
	Teardown() error
}

// Note that this interface, while descrabing what is necessary to be a game, does not describe the specific
// functionality which you should expect of a regular game. Instead, look at game_test.go to understand that
// a little better. There are cases, such as one or zero players and a game that is ongoing, which are not
// valid and stopped by specific play implementions or the servers which run them. These are either tested on
// server tests or not tested at all since they are not desired configurations. Alternatively, they may be tested
// of specific implementations.

type GameInitArgs struct {
	// Human-Identifyiers
	Name     *string
	JoinCode *string

	// Game Settings
	Public        bool 
	MaxPlayers    uint64
	Stakes        uint64
	StartingChips uint64
	Mode          uint64

	// Renew Information (if you keep chips you must keep players)
	KeepPlayers bool
	KeepChips   bool
}
