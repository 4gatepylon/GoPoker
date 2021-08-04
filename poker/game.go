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
	PSTATUS_ADMIN byte = 1 << iota
	PSTATUS_PLAYING
)

// Game Status
const (
	GSTATUS_PLAYING byte = 1 << iota
	GSTATUS_PRIVATE
)

// Game Mode
const (
	GMODE_CONST_STAKES uint64 = 1 << iota
)

// Defaults for standard games
const (
	DEFAULT_STAKES                 = 1000               // Default big blind is 1000 chips
	DEFAULT_MAX_PLAYERS            = 6                  // By default allow six players maximum
	DEFAULT_MODE                   = GMODE_CONST_STAKES // Standard games have constant stakes.
	DEFAULT_STATUS                 = 0                  // The default status is the null status (not started)
	DEFAULT_STAKES_HAND_MULTIPLIER = 100                // DEFAULT_STAKES * ..._MULTIPLIER = default starting hand
)

// A GameLike should be able to manipulate CardLikes accordingly. The string method
// will be desired to communiate with players.
type CardLike interface {
	String() string
}

type PlayerInfo struct {
	// Human-Identifiers
	Name  *string
	Chips uint64
	Pot   uint64

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
	AddPlayer(*string) (*string, bool, error) // (prospective player name) => (player name, joined, error)
	Players() []*PlayerInfo                   // () => (an informative list of players in order of play)
	Stakes() uint64                           // () => (value of big blind in chips)
	Middle() *[5]CardLike                     // () => (array of cards in the middle)

	// Game Status
	ChangeGameName(*string, *string) (bool, error) // (name changer, desired name) => (changed name, error)
	Play(*string) (bool, error)                    // (play requester) => (played, error)
	Pause(*string) (bool, error)                   // (pause requester) => (paused, error)
	MakePrivate(*string) (bool, error)             // (make private requester) => (made private, error)
	MakePublic(*string) (bool, error)              // (make public requester) => (made public, error)
	Playing() bool                                 // () => (game is playing)
	Private() bool                                 // () => (game is private)

	// Game Flow
	Move(uint64, uint64, *string) (bool, error)               // (move, chips: optional, mover) => (moved, error)
	ChangePlayerName(*string, *string) (*string, bool, error) // (namer, player) => (new name, renamed, error)

	// Game Flow Control Plane
	Increment() (bool, error)  // () => (incremented, error)
	Resolve() (*string, error) // () => (winners' informative message, error)
	NewRound() error           // () => (error)
	Renew() error              // () => (error)
}

type GameInitArgs struct {
	// Human-Identifyiers
	Name     *string
	JoinCode *string

	// Game Settings
	Public        bool
	MaxPlayers    uint64
	StartingChips uint64
	Stakes        uint64
	Mode          uint64

	// Renew Information
	KeepPlayers bool
	KeepChips   bool
}
