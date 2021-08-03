package poker

import (
	"fmt"
	"github.com/4gatepylon/GoPoker/utils"
)

type Card CardSet

const (
	BROUND_PREFLOP uint64 = iota + 1
	BROUND_FLOP
	BROUND_TURN
	BROUND_RIVER
)

const (
	MTYPE_CHECK uint64 = 1 << iota // Takes presedence over other moves
	MTYPE_FOLD                     // Second highest precedence
	MTYPE_CALL                     // Third highest precedence 
	MTYPE_CALL_ANY                 // Call anything (willing to go all in)
	MTYPE_BET                      // Bet some amount of money; you can check AND bet
	MTYPE_SITOUT_NEXT_ROUND        // Orthogonal from the rest: simply don't lay the next round (toggle pause)
)

const (
	PSTATUS_ADMIN   byte = 1 << iota // Admins control the chips and playerbase
	PSTATUS_PLAYING                  // Players who aren't playing can't bet (etc)
)

const (
	GSTATUS_PLAYING byte = 1 << iota // Games can be in play or in pause
	GSTATUS_PRIVATE                  // and it can be private or public (private requires JoinCode != nil)
)

const (
	GMODE_CONST_STAKES uint64 = 1 << iota // Game mode in which the stakes are constant
	GMODE_RAKE                            // Whether the house will be raking this game
)

const (
	DEFAULT_STAKES      = 1000
	DEFAULT_MAX_PLAYERS = 0
	DEFAULT_RAKE_AMOUNT = 0
	DEFAULT_MODE        = GMODE_CONST_STAKES
	DEFAULT_STATUS      = 0
)


// TODO finalize this interface
// TODO thread safety, rework current functions
// TODO finish implementing new functions
// TODO unit-testing (at least in a single-threaded environment)

type GameLike interface {
	// Player control (these will be used by admins/mods)
	AddPlayer()        // Add a player to the game (mod required)
	KickPlayer()       // Remove a player from the game (mod required)
	ModPlayer()        // Turn a player into a mod (admin) (mod required)
	

	// Player control plane (this will be used by software)
	Players() // Who is playing the game: return in order with chips
	Stakes()  // What are the stakes? (bb)

	// Game status control plane (this will be used by admins/mods)
	ChangeGameName()
	Play()           // Toggle pause OFF
	Pause()          // Toggle pause ON
	MakePrivate()    // Toggle private ON (join code will be necessary)
	MakePublic()     // Toggle private OFF (join code will be ignored)
	Private()        // Read whether this game is private or not

	// Game flow (this will be used by everyone
	Move(uint64, uint64, uint64) // Used to make bets, checks, and so on
	ChangePlayerName()           // Used to change a name: admins can change anyone's, and anyone else can change their own

	// Game flow control plane (this will be used by software)
	Increment() // Increment the betting round (preflop => flop => turn => river)
	Resolve()   // Resolve the winners
	NewRound()  // Start the next round

	// NOTE: resolve should be able to handle different pots and all ins (incl. forced all ins)
	// Moreover, a game has to remember player order, it has to be able to deal with betting rounds within 
	// game rounds within the game itself, etc...
}

// Internally we store players as their own struct but player structs are not "able to do anything"
// on their own. They are just used to keep track of cards, names, chips, etc...
type Player struct {
	Id   uint64  // Players have unique identifiers for the system
	Name *string // and display names for people

	Hand   [2]Card // A player has a two-card hand
	Chips  uint64  // and a positive number of chips
	Status byte    // and a status (i.e. is this player an admin? is he playing?)
	GameId uint64  // Each player is in a game or in the zero game id, which is lobby
}

type Game struct {
	JoinCode *string // A join code is effectively a password to join a game
	Name     *string // A descriptive name for the game (shows up if you query for games on a server)

	Id         uint64 // Games are uniquely identified for the system
	MaxPlayers uint64 // Games must cap the number of players; zero means uncapped

	Middle  [5]Card            // Cards in the middle (self explanatory)
	Players map[uint64]*Player // Up to the MaxPlayers number of players (recommended is six)
	Status  byte               // Private | Public, Playing | Paused, etc...

	Mode          uint64 // The game mode (i.e. constant stakes, rake | no rake)
	Stakes        uint64 // The Value of big blind (3x little blind)
	RakeAmount    uint64 // The amount of rake (what the house takes every round, additional to blinds, from all players)
	StartingChips uint64 // The amount of chips to give to new players when they join (default will be 10x bb)
}

// Create a new game
type GameInitArgs struct {
	Name     *string // The name of the game to display to queries in the lobby
	JoinCode *string // The password needed to join the game

	Public        bool   // Whether anybody at all can join the game or not (will need join code)
	MaxPlayers    uint64 // The maximum number of players allowed in the game at a time (a game is a table)
	StartingChips uint64 // The number of chips to assign everybody to start with
	Stakes        uint64 // The value of big blind (3x little blind)
	Mode          uint64 // The game mode (default is zero, which will be rerouted to GMODE_CONST_STAKES)
	RakeAmount    uint64 // How much the house will rake

	KeepPlayers bool // Used for renew: if players already exist, whether to keep them or kick them
	KeepChips   bool // Used for renew: if players already exist, whether to keep the chips or not
}

func New(creator string, args *GameInitArgs) (*Game, error) {
	joinCode := args.JoinCode
	if joinCode == nil {
		jc := fmt.Sprintf("%s-%s", utils.RandVerbAdv(nil), utils.RandString(3))
		joinCode = &jc
	}

	name := args.Name
	if name == nil {
		n := fmt.Sprintf("%s's game", creator)
		name = &n
	}
	maxPlayers := args.MaxPlayers // Recall zero is unlimited
	stakes := args.Stakes
	if stakes == 0 {
		stakes = DEFAULT_STAKES
	}
	rakeAmount := args.RakeAmount // Not supported
	startingChips := args.StartingChips
	if startingChips == 0 {
		startingChips = 10 * stakes
	}

	g := &Game{
		JoinCode:      joinCode,                                             // ...
		Name:          name,                                                 // ...
		Middle:        [5]Card{NoCards, NoCards, NoCards, NoCards, NoCards}, // Empty until we start the game
		Id:            utils.RandInt64(),                                    // A random id
		MaxPlayers:    maxPlayers,                                           // As above (zero is infinite)
		Players:       nil,                                                  // Initialized lazily as we add players
		Status:        DEFAULT_STATUS,                                       // Status 0 simply is a negation of all statuses
		Mode:          DEFAULT_MODE,                                         // Rake not yet supported
		Stakes:        stakes,                                               // ...
		RakeAmount:    rakeAmount,                                           // Not yet supported
		StartingChips: startingChips,                                        // ...
	}

	g.AddPlayer(&creator)
	g.ModPlayer(&creator, PSTATUS_ADMIN)

	return g, nil
}

// REQUIRES ADMIN
func (game *Game) ChangeName(name *string) {
	if name == nil {
		n := fmt.Sprintf("anonymous%s's game", utils.RandString(3))
		name = &n
	}
	game.Name = name
}

// REQUIRES ADMIN
func (game *Game) changeJoinCode(joinCode *string) {
	if joinCode == nil {
		jc := fmt.Sprintf("%s-%s", utils.RandVerbAdv(nil), utils.RandString(3))
		joinCode = &jc
	}
	game.JoinCode = joinCode
}

// REQUIRES ADMIN
func (game *Game) changeMaxPlayers(maxPlayers uint64) {
	game.MaxPlayers = maxPlayers
}

// REQUIRES ADMIN
func (game *Game) changeStartingChips(startingChips uint64) {
	game.StartingChips = startingChips
}

// REQUIRES ADMIN
func (game *Game) TogglePause() {
	game.Status = (game.Status ^ GSTATUS_PLAYING)
}
func (game *Game) Pause() {
	game.Status = game.Status | GSTATUS_PLAYING
}
func (game *Game) Play() {
	game.Status = game.Status & ^GSTATUS_PLAYING
}

// REQUIRES ADMIN
func (game *Game) TogglePrivate() {
	game.Status = (game.Status ^ GSTATUS_PRIVATE)
}
func (game *Game) MakePrivate() {
	game.Status = game.Status | GSTATUS_PRIVATE
}
func (game *Game) MakePublic() {
	game.Status = game.Status & ^GSTATUS_PRIVATE
}

// Check whether a player can be added with a joinCode
func (game *Game) ValidJoinCode(joinCode *string) bool {
	return (game.Status&GSTATUS_PRIVATE > 0) || (joinCode != nil && game.JoinCode != nil && *joinCode == *game.JoinCode)
}

func randPlayerName() string {
	adjAnimal := utils.RandAdjAnimal(nil) // default separator is "-"
	// >= 7 * log(64) bits of randomness from the string + >= log(16)^2 from the adj/animal
	// ~ >= 3 * 8 + 4 + 4 = 8 * 8 = 32... which should be enough for regular games with <= usually 10 ppl
	return fmt.Sprintf("%s-%s", adjAnimal, utils.RandString(3))
}

// Add a player; a nil name will create a random name
func (game *Game) AddPlayer(namep *string) (*string, error) {
	name := randPlayerName()
	if namep != nil {
		name = *namep
	}

	p := &Player{
		Hand:   [2]Card{NoCards, NoCards},
		Chips:  game.StartingChips,
		Status: PSTATUS_PLAYING,
		Id:     utils.RandInt64(),
		Name:   &name,
		GameId: game.Id,
	}
	if game.Players == nil {
		game.Players = make(map[uint64]*Player, 1)
		game.Players[p.Id] = p
		return p.Name, nil
	}
	for _, player := range game.Players {
		if p.Id == player.Id {
			return nil, fmt.Errorf("Player with id %d exists\n", p.Id)
		}
		if p.Name == player.Name {
			return nil, fmt.Errorf("Player with name %s exists\n", name)
		}
	}
	game.Players[p.Id] = p
	return p.Name, nil
}

func getPlayer(g *Game, id interface{}) (*Player, error) {
	switch id.(type) {
	case uint64:
		intId, _ := id.(uint64) // Cannot fail inside this case
		p, ok := g.Players[intId]
		if !ok {
			return nil, nil
		}
		return p, nil
	case *string:
		pid, _ := id.(*string) // Similarly, cannot fail inside the case
		if pid == nil {
			return nil, fmt.Errorf("Tried to get nil-named player\n")
		}
		for _, player := range g.Players {
			if *pid == *player.Name {
				return player, nil
			}
		}
	default:
		return nil, fmt.Errorf("Passed in type %T but expected uint64 (id) or *string (name)\n in getPlayer", id)
	}
	return nil, fmt.Errorf("Unreachable code in getPlayer helper: missed type switch\n")
}

// REQUIRES ADMIN
func (game *Game) KickPlayer(id interface{}) error {
	player, err := getPlayer(game, id)
	if err != nil {
		return err
	}
	if player == nil {
		return fmt.Errorf("Tried to kick nil player\n")
	}

	delete(game.Players, player.Id)
	return nil
}

// Change the status of a player by id or name
// REQUIRES ADMIN
func (game *Game) ModPlayer(id interface{}, status byte) error {
	player, err := getPlayer(game, id)
	if err != nil {
		return err
	}
	if player == nil {
		return fmt.Errorf("Tried to mod nil player\n")
	}

	player.Status |= status
	return nil
}

// Check whether a player is able to give chips, change name, etc...
func (game *Game) IsAdmin(id interface{}) (bool, error) {
	player, err := getPlayer(game, id)
	if err != nil {
		return false, err
	}
	if player == nil {
		return false, fmt.Errorf("Tried to query admin of nil player\n")
	}

	return player.Status&PSTATUS_ADMIN > 0, nil
}

// Give a player (by id or name) chips
// REQUIRES ADMIN
func (game *Game) GiveChips(id interface{}, chips uint64) error {
	player, err := getPlayer(game, id)
	if err != nil {
		return err
	}
	if player == nil {
		return fmt.Errorf("Tried to mod nil player\n")
	}

	player.Chips += chips
	return nil
}

// Attempt to make a move with some chips; chips are ignored for checks and folds
func (game *Game) Move(move uint64, pid uint64, chips uint64) {

}

// Resolve the winners from the current middle and those playing
func (game *Game) Resolve() {

}

// Start a new round
// Should only be possible in the river when there are no chips in the middle
// and everyone has checked or called
func (game *Game) NewRound(move uint64) {

}

// Increment the betting round (i.e. preflop => flop => turn => river)
// Should only be possible after everyone in the game has bet
func (game *Game) Increment() {

}