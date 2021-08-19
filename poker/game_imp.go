package poker

import (
	"fmt"
	"os"
	"log"
	"io/ioutil"
	"path/filepath"
	"strings"
	"encoding/json"
	"github.com/4gatepylon/GoPoker/utils"
)

// FIXME thread safety, logging, reworking (data structure too/efficiency/speed), make sure it works..., simplify!

const maxUint64 uint64 = 0xffffffffffffffff

type Card CardSet

func (c Card) String() string {
	return CardSetToString(CardSet(c))
}

// Internally we store players as their own struct but player structs are not "able to do anything"
// on their own. They are just used to keep track of cards, names, chips, etc...
type Player struct {
	Id   uint64    // Players have unique identifiers for the system
	Name *string   // and display names for people

	Hand   [2]Card // A player has a two-card hand
	Chips  uint64  // and a positive number of chips
	Bet    uint64  // Number of chips in play that are bet
	Pot    uint64  // Number of chips in the pot not being bet
	Status uint64  // and a status (i.e. is this player an admin? is he playing?)
	GameId uint64  // Each player is in a game or in the zero game id, which is lobby
}

type Pot struct {
	Chips uint64
	Players []uint64
}

type Game struct {
	Id         uint64     // Games are uniquely identified for the system

	joinCode      *string // A join code is effectively a password to join a game
	name          *string // A descriptive name for the game (shows up if you query for games on a server)
	maxPlayers    uint64  // Games must cap the number of players; zero means uncapped

	// Game control
	middle        [5]Card              // Cards in the middle (self explanatory)
	players       map[uint64]*Player   // Up to the MaxPlayers number of players (recommended is six)
	status        uint64               // Private | Public, Playing | Paused, etc...
	pots          []Pot
	bettingRound  uint64
	roundNum      uint64

	mode          uint64 // The game mode (i.e. constant stakes)
	stakes        uint64 // The Value of big blind (3x little blind)
	startingChips uint64 // The amount of chips to give to new players when they join (default will be 10x bb)
	
	// Maintenance
	gameDir       *string
	errorLog      *os.File
	roundLog      *os.File
	gameInit      *os.File
	errorLogger   *log.Logger
	roundLogger   *log.Logger
}

const errorLogName = "error.log"
const roundLogName = "round.log"
const gameInitName = "info.json"

// Human readable encoding (json) for gameInit file
type gameInitJson struct {
	Id            string   `json: game-id`
	Name          string   `json: game-name`
	JoinCode      string   `json: game-join-code`
	Status        string   `json: game-status`
	MaxPlayers    uint64   `json: max-players`
	Stakes        uint64   `json: big-blind`
	StartingChips uint64   `json: starting-chips`
	Mode          string   `json: game-mode`
	KeepPlayers   bool    `json: keep-players-on-renew`
	KeepChips     bool    `json: keep-chips-on-renew`
}

func gameMode2Str(mode uint64) (string, error) {
	if mode == GMODE_CONST_STAKES {
		return "CONST_STAKES", nil
	}
	return "", fmt.Errorf("Invalid game mode: %d", mode)
}

func str2GameMode(mode string) (uint64, error) {
	if mode == "CONST_STAKES" {
		return GMODE_CONST_STAKES, nil
	}
	return 0, fmt.Errorf("Invalid game mode (str): %d", mode)
}

func gameStatus2Str(status uint64) (string, error) {
	statuses := make([]string, 0, 2)
	if status | GSTATUS_PLAYING > 0 {
		statuses = append(statuses, "PLAYING")
	} else {
		statuses = append(statuses, "PAUSED")
	}
	if status | GSTATUS_PRIVATE > 0 {
		statuses = append(statuses, "PRIVATE")
	} else {
		statuses = append(statuses, "PUBLIC")
	}
	if len(statuses) == 0 {
		return "", fmt.Errorf("Got no statuses")
	}
	return strings.Join(statuses, " "), nil
}

func str2GameStatus(statusStr string) (uint64, error) {
	var status uint64
	statuses := strings.Split(statusStr, " ")
	if len(statuses) == 0 || len(statuses) > 2 {
		return 0, fmt.Errorf("Got unknown number of statuses: %d", len(statuses))
	}
	// Order matters here
	valid := [][2]string{[2]string{"PAUSED", "PLAYING"}, [2]string{"PUBLIC", "PRIVATE"}}
	for i, v := range valid {
		if statuses[i] == v[0] {
			// Check game.go for this
			status |= (1 << i)
		} else if statuses[i] != v[1] {
			return 0, fmt.Errorf("Found status %s, but should have been %s or %s", statuses[i], v[0], v[1])
		}
	}
	return status, nil
}

func New(creator *string, args *GameInitArgs) (*string, GameLike, error) {
	// Initialize all the settings using defaults if they don't provide any
	joinCode := args.JoinCode
	if joinCode == nil {
		jc := fmt.Sprintf("%s-%s", utils.RandVerbAdv(nil), utils.RandString(3))
		joinCode = &jc
	}

	if creator == nil {
		rn := randPlayerName()
		creator = &rn
	}

	id := utils.RandInt64()
	name := args.Name
	if name == nil {
		n := fmt.Sprintf("%ss-game", *creator)
		name = &n
	}

	status := DEFAULT_STATUS
	if !args.Public {
		status |= GSTATUS_PRIVATE
	}

	maxPlayers := args.MaxPlayers
	if maxPlayers == 0 {
		maxPlayers = DEFAULT_MAX_PLAYERS
	}

	stakes := args.Stakes
	if stakes == 0 {
		stakes = DEFAULT_STAKES
	}

	startingChips := args.StartingChips
	if startingChips == 0 {
		startingChips = 10 * stakes
	}

	if args.Mode != 0 && args.Mode != DEFAULT_MODE {
		return nil, nil, fmt.Errorf("Tried to create game with mode %d, but only mode %d supported", args.Mode, DEFAULT_MODE)
	}

	// Create directory with game information
	gameDir, err := ioutil.TempDir("", fmt.Sprintf("%s-*", *name))
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to make tempDir: `%v`", err)
	}

	// We use a standard naming scheme outlined below
	// 0744 is 7 => all perms for this program, 4 => only read for others
	errorLog, err  := os.OpenFile(filepath.Join(gameDir, errorLogName), os.O_RDWR|os.O_CREATE, 0744)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to open errorLog: `%v`", err)
	}
	errorLogger := log.New(errorLog, "", log.Lshortfile | log.Ltime | log.LUTC)
	roundLog, err := os.OpenFile(filepath.Join(gameDir, roundLogName), os.O_RDWR|os.O_CREATE, 0744)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to open roundLog: `%v`", err)
	}
	roundLogger := log.New(roundLog, "", log.Ltime | log.LUTC)

	// Store the game initialization parameters in case we crash in the middle of the game
	gameInit, err := os.OpenFile(filepath.Join(gameDir, gameInitName), os.O_RDWR|os.O_CREATE, 0776)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to open gameInfo: `%v`", err)
	}
	// We store the generated values like joinCode and name
	gs, err := gameStatus2Str(status)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to convert status %d to string: `%v`", status, err)
	}
	gm, err := gameMode2Str(DEFAULT_MODE)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to convert mode %d to string: `%v`", DEFAULT_MODE, err)
	}
	m, err := json.MarshalIndent(gameInitJson{
		Id:            fmt.Sprintf("%d", id),
		Name:          *name,
		JoinCode:      *joinCode,
		Status:        gs,
		MaxPlayers:    maxPlayers,
		Stakes:        stakes,
		StartingChips: startingChips,
		Mode:          gm,
		KeepPlayers:   args.KeepPlayers,
		KeepChips:     args.KeepChips,
	}, "", "  ")
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to marshal game init args: `%v`", err)
	}
	// Recall that in go if it fails to write ALL the bytes it will error out
	_, err = gameInit.Write(m)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to store game init args: `%v`", err)
	}
	err = gameInit.Close()
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to close game info: `%v`", err)
	}

	// Create the game object in a state that is NOT yet playing
	g := &Game{
		joinCode:      joinCode,
		name:          name, 
		middle:        [5]Card{NoCards, NoCards, NoCards, NoCards, NoCards},
		Id:            id, // A random id
		maxPlayers:    maxPlayers,        // As above (zero is infinite)
		players:       nil,               // Initialized lazily as we add players
		status:        status,            // Status 0 simply is a negation of all statuses
		mode:          DEFAULT_MODE,      // Rake not yet supported
		stakes:        stakes,            // ...
		startingChips: startingChips,     // ...
		gameDir:       &gameDir,
		errorLog:      errorLog,
		roundLog:      roundLog,
		gameInit:      gameInit,
		errorLogger:   errorLogger,
		roundLogger:   roundLogger,
	}

	// Add the creator as an admin
	_, added, err := g.AddPlayer(creator, joinCode)
	if !added || err != nil {
		return nil, nil, fmt.Errorf("Failed to add (added = %v) creator `%s`: `%v`", added, *creator, err)
	}
	modded, err := g.ModPlayer(nil, creator, PSTATUS_ADMIN)
	if !modded || err != nil {
		return nil, nil, fmt.Errorf("Failed to mod (modded=%v) creator `%s` to admin status: `%v`", modded, *creator, err)
	}

	return creator, g, nil
}

func (g *Game) Renew() error {
	// TODO
	return fmt.Errorf("Not implemented!")
}

func (g *Game) Teardown() error {
	err := g.errorLog.Close()
	if err != nil {
		return fmt.Errorf("Failed to close error log: `%v`", err)
	}
	err = g.roundLog.Close()
	if err != nil {
		return fmt.Errorf("Failed to close roundLog: `%v`", err)
	}
	return os.RemoveAll(*g.gameDir)
}

func randPlayerName() string {
	adjAnimal := utils.RandAdjAnimal(nil) // default separator is "-"
	return fmt.Sprintf("%s-%s", adjAnimal, utils.RandString(3))
}

func (g *Game) getPlayer(name *string) (*Player, bool) {
	for _, p := range g.players {
		if *p.Name == *name {
			return p, true
		}
	}
	return nil, false
}

func (g *Game) isAdmin(name *string) (bool, error) {
	p, found := g.getPlayer(name)
	if !found {
		return false, fmt.Errorf("Could not find player %s", *name)
	}
	return (p.Status & PSTATUS_ADMIN) > 0, nil
}

func (g *Game) onlyExecuteIfIsAdmin(admin *string, f func() (bool, error)) (bool, error) {
	mod, err := g.isAdmin(admin)
	if err != nil {
		return false, fmt.Errorf("Failed to check if player %s was admin: %v", *admin, err)
	}
	if !mod {
		return false, nil
	}
	return f()
}

// Add a player; a nil name will create a random name
func (g *Game) AddPlayer(name *string, joinCode *string) (*string, bool, error) {
	if !g.Private() && g.joinCode == nil {
		return nil, false, fmt.Errorf("Trying to add to a nil joincode game that is private")
	}
	if g.Private() && (joinCode != nil || *joinCode != *g.joinCode) {
		return nil, false, nil
	}
	if name == nil {
		n := randPlayerName()
		name = &n
	}
	p := &Player{
		Hand:   [2]Card{NoCards, NoCards},
		Chips:  g.startingChips,
		Status: 0, // Not yet playing
		Id:     utils.RandInt64(),
		Name:   name,
		GameId: g.Id,
	}
	if g.players != nil {
		for _, player := range g.players {
			if p.Id == player.Id {
				return nil, false, fmt.Errorf("Player with id %d exists\n", p.Id)
			}
			if p.Name == player.Name {
				return nil, false, fmt.Errorf("Player with name %s exists\n", name)
			}
		}
	} else {
		g.players = make(map[uint64]*Player, 1)
	}
	g.players[p.Id] = p
	return p.Name, true, nil
}

func (g *Game) KickPlayer(kicker *string, kicked *string) (bool, error) {
	return g.onlyExecuteIfIsAdmin(kicker, func() (bool, error) {
		rec, found := g.getPlayer(kicked)
		if !found {
			return false, fmt.Errorf("Tried to kick nonexistent player %s", *kicked)
		}
		// FIXME: add some checking for whether the game is in play or not (etc)
		delete(g.players, rec.Id)
		return true, nil
	})
}

// Change the status of a player by id or name
func (g *Game) ModPlayer(modder *string, modded *string, mod uint64) (bool, error){
	return g.onlyExecuteIfIsAdmin(modder, func() (bool, error) {
		rec, found := g.getPlayer(modded)
		if !found {
			return false, fmt.Errorf("Did not find player %s to mod", *modded)
		}
		// We may want to change this to add and remove permissions later
		rec.Status = mod
		return true, nil
	})
}

func (g *Game) ChangePlayerName(changer *string, name *string, newName *string) (*string, bool, error) {
	if newName == nil {
		n := randPlayerName()
		newName = &n
	}
	changed, err := g.onlyExecuteIfIsAdmin(changer, func() (bool, error) {
		rec, found := g.getPlayer(name)
		if !found {
			return false, fmt.Errorf("Did not find player %s to change name for", *name)
		}
		rec.Name = newName
		return true, nil
	})
	return newName, changed, err
}

func (g *Game) ChangeGameName(changer *string, name *string) (bool, error) {
	return g.onlyExecuteIfIsAdmin(changer, func() (bool, error) {
		if name == nil {
			n := fmt.Sprintf("%s-game-%s", changer)
			name = &n
		}
		g.name = name
		return true, nil
	})
}

func (g *Game) GiveChips(giver *string, receiver *string, chips uint64) (bool, error) {
	return g.onlyExecuteIfIsAdmin(giver, func() (bool, error) {
		rec, found := g.getPlayer(receiver)
		if !found {
			return false, fmt.Errorf("Did not find player %s to send the chips to", *receiver)
		}
		if maxUint64 - chips < rec.Chips {
			return false, fmt.Errorf("Chips would overflow storage medium")
		}
		rec.Chips += chips
		return true, nil
	})
}

func (g *Game) Stakes() uint64 {
	return g.stakes
}

func (g *Game) Pots() []uint64 {
	pots := make([]uint64, len(g.pots))
	for i, _ := range pots {
		pots[i] = g.pots[i].Chips
	}
	return pots
}

func (g *Game) Playing() bool {
	return g.status & GSTATUS_PLAYING > 0
}

func (g *Game) Private() bool {
	return g.status & GSTATUS_PRIVATE > 0
}

func (g *Game) Players() []*PlayerInfo {
	players := make([]*PlayerInfo, 0, len(g.players))
	for _, p := range g.players {
		var c [2]CardLike
		for i, _ := range p.Hand {
			c[i] = p.Hand[i]
		}
		players = append(players, &PlayerInfo{
			Name:  *p.Name,
			Chips: p.Chips,
			Bet:   p.Bet,
			Id:    p.Id,
			Cards: c,
			Mod:   (p.Status & PSTATUS_ADMIN) > 0,
		})
	}
	return players
}

func (g *Game) Middle() *[5]CardLike {
	var cp [5]CardLike
	for i, _ := range g.middle {
		cp[i] = g.middle[i]
	}
	return &cp
}

func (g *Game) Pause(pauser *string) (bool, error) {
	return g.onlyExecuteIfIsAdmin(pauser, func() (bool, error) {
		g.status = g.status & ^GSTATUS_PLAYING
		return true, nil
	})
}
func (g *Game) Play(player *string) (bool, error) {
	return g.onlyExecuteIfIsAdmin(player, func() (bool, error) {
		g.status = g.status | GSTATUS_PLAYING
		return true, nil
	})
}

func (g *Game)  MakePrivate(privater *string) (bool, error) {
	return g.onlyExecuteIfIsAdmin(privater, func() (bool, error) {
		g.status = g.status | GSTATUS_PRIVATE
		return true, nil
	})
}
func (g *Game) MakePublic(publicer *string) (bool, error) {
	return g.onlyExecuteIfIsAdmin(publicer, func() (bool, error) {
		g.status = g.status & ^GSTATUS_PRIVATE
		return true, nil
	})
}

//////////////////////////////////////////////////////////////////// Game flow functionality

// Attempt to make a move with some chips; chips are ignored for checks and folds
func (g *Game) Move(move uint64, chips uint64, mover *string) (bool, error) {
	// TODO
	return false, nil
}

func (g *Game) Increment() (bool, error) {
	return false, nil // TODO
}

// Resolve the winners from the current middle and those playing
func (g *Game) Resolve() (*string, error) {
	// NOTE: resolve should be able to handle different pots and all ins (incl. forced all ins)
	// Moreover, a game has to remember player order, it has to be able to deal with betting rounds within
	// game rounds within the game itself, etc...

	// TODO
	return nil, nil
}

// Start a new round
// Should only be possible in the river when there are no chips in the middle
// and everyone has checked or called
func (g *Game) NewRound() error {
	// TODO
	return nil
}