package poker

import (
	"fmt"
	"github.com/4gatepylon/GoPoker/utils"
)

// TODO thread safety, rework current functions
// TODO finish implementing new functions
// TODO unit-testing (at least in a single-threaded environment)

type Card struct {
	c CardSet
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

	Mode          uint64 // The game mode (i.e. constant stakes)
	Stakes        uint64 // The Value of big blind (3x little blind)
	StartingChips uint64 // The amount of chips to give to new players when they join (default will be 10x bb)
}

func New(creator *string, args *GameInitArgs) (GameLike, error) {
	joinCode := args.JoinCode
	if joinCode == nil {
		jc := fmt.Sprintf("%s-%s", utils.RandVerbAdv(nil), utils.RandString(3))
		joinCode = &jc
	}

	name := args.Name
	if name == nil {
		n := fmt.Sprintf("%s's game", *creator)
		name = &n
	}
	maxPlayers := args.MaxPlayers // Recall zero is unlimited
	stakes := args.Stakes
	if stakes == 0 {
		stakes = DEFAULT_STAKES
	}
	startingChips := args.StartingChips
	if startingChips == 0 {
		startingChips = 10 * stakes
	}

	g := &Game{
		JoinCode: joinCode, // ...
		Name:     name,     // ...
		// TODO
		// Middle:        [5]Card{NoCards, NoCards, NoCards, NoCards, NoCards}, // Empty until we start the game
		Id:            utils.RandInt64(), // A random id
		MaxPlayers:    maxPlayers,        // As above (zero is infinite)
		Players:       nil,               // Initialized lazily as we add players
		Status:        DEFAULT_STATUS,    // Status 0 simply is a negation of all statuses
		Mode:          DEFAULT_MODE,      // Rake not yet supported
		Stakes:        stakes,            // ...
		StartingChips: startingChips,     // ...
	}

	g.AddPlayer(creator)
	g.ModPlayer(creator, PSTATUS_ADMIN)

	return nil, nil
}

// REQUIRES ADMIN
func (game *Game) ChangeGameName(name *string) {
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
		// Hand:   [2]Card{NoCards, NoCards},
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
	// TODO
}

// Resolve the winners from the current middle and those playing
func (game *Game) Resolve() {
	// NOTE: resolve should be able to handle different pots and all ins (incl. forced all ins)
	// Moreover, a game has to remember player order, it has to be able to deal with betting rounds within
	// game rounds within the game itself, etc...

	// TODO
}

// Start a new round
// Should only be possible in the river when there are no chips in the middle
// and everyone has checked or called
func (game *Game) NewRound(move uint64) {
	// TODO
}

// Increment the betting round (i.e. preflop => flop => turn => river)
// Should only be possible after everyone in the game has bet
func (game *Game) Increment() {
	// TODO
}
