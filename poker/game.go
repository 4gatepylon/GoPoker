package poker

import (
	"github.com/google/uuid"
	"fmt"
	"github.com/4gatepylon/GoPoker/utils"
	"time"
	"math/rand"
)

const maxUint64 uint64 = 0xffffffffffffffff

// Above this number, it's going to be hard to write a GUI of any kind
// Note that above 52 - 2 * n - 5 => 47 - 2n; n >= 23 players we need more decks
const foreverMaxPlayers int = 12

// Each player gets a max bet time
const bettingDeadline time.Duration = 45 * time.Second

// Each player will get some time to see the outcome of the round
const newRoundDeadline time.Duration = 15 * time.Second

// NOTE: a lot of the timing will be done client-side
// and there will be a need to synchronize clocks... it's optional to have
// the deadline entirely on the client and then forward it here... I'll look into it
// for now these consts will not be used

// Hand types
const (
	_ = iota
	HTYPE_HIGHCARD
	HTYPE_PAIR
	HTYPE_TWO_PAIR
	HTYPE_TRIP
	HTYPE_STRAIGHT
	HTYPE_FLUSH
	HTYPE_FULLHOUSE
	HTYPE_STRAIGHT_FLUSH
	HTYPE_QUAD
	HTYPE_ROYAL_FLUSH
)

func HandTypeStr2Int(htype string) (int, error) {
	switch htype {
	case "HTYPE_HIGHCARD":
		return HTYPE_HIGHCARD, nil
	case "HTYPE_PAIR":
		return HTYPE_PAIR, nil
	case "HTYPE_TWO_PAIR":
		return HTYPE_TWO_PAIR, nil
	case "HTYPE_TRIP":
		return HTYPE_TRIP, nil
	case "HTYPE_STRAIGHT":
		return HTYPE_STRAIGHT, nil
	case "HTYPE_FLUSH":
		return HTYPE_FLUSH, nil
	case "HTYPE_FULLHOUSE":
		return HTYPE_FULLHOUSE, nil
	case "HTYPE_STRAIGHT_FLUSH":
		return HTYPE_STRAIGHT_FLUSH, nil
	case "HTYPE_QUAD":
		return HTYPE_QUAD, nil
	case "HTYPE_ROYAL_FLUSH":
		return HTYPE_ROYAL_FLUSH, nil
	default:
		return 0, fmt.Errorf("Invalid hand type")
	}
}

func HandTypeInt2String(htype int) (string, error) {
	switch htype {
	case HTYPE_HIGHCARD:
		return "HTYPE_HIGHCARD", nil
	case HTYPE_PAIR:
		return "HTYPE_PAIR", nil
	case HTYPE_TWO_PAIR:
		return "HTYPE_TWO_PAIR", nil
	case HTYPE_TRIP:
		return "HTYPE_TRIP", nil
	case HTYPE_STRAIGHT:
		return "HTYPE_STRAIGHT", nil
	case HTYPE_FLUSH:
		return "HTYPE_FLUSH", nil
	case HTYPE_FULLHOUSE:
		return "HTYPE_FULLHOUSE", nil
	case HTYPE_STRAIGHT_FLUSH:
		return "HTYPE_STRAIGHT_FLUSH", nil
	case HTYPE_QUAD:
		return "HTYPE_QUAD", nil
	case HTYPE_ROYAL_FLUSH:
		return "HTYPE_ROYAL_FLUSH", nil
	default:
		return "", fmt.Errorf("Invalid hand type")
	}
}

// Betting rounds
const (
	BROUND_NONE = iota
	BROUND_PREFLOP
	BROUND_FLOP
	BROUND_TURN
	BROUND_RIVER
)

func BettingRoundStr2Int(bround string) (int, error) {
	switch bround {
	case "BROUND_NONE":
		return BROUND_NONE, nil
	case "BROUND_PREFLOP":
		return BROUND_PREFLOP, nil
	case "BROUND_FLOP":
		return BROUND_FLOP, nil
	case "BROUND_TURN":
		return BROUND_TURN, nil
	case "BROUND_RIVER":
		return BROUND_RIVER, nil
	default:
		return 0, fmt.Errorf("Unknown betting round")
	}
}

func BettingRoundInt2Str(bround int) (string, error) {
	switch bround {
	case BROUND_NONE:
		return "BROUND_NONE", nil
	case BROUND_PREFLOP:
		return "BROUND_PREFLOP", nil
	case BROUND_FLOP:
		return "BROUND_FLOP", nil
	case BROUND_TURN:
		return "BROUND_TURN", nil
	case BROUND_RIVER:
		return "BROUND_RIVER", nil
	default:
		return "", fmt.Errorf("Unknown betting round")
	}
}

// Move Types
const (
	MTYPE_CHECK = 1 << iota
	MTYPE_FOLD
	MTYPE_CALL
	MTYPE_CALL_ANY
	MTYPE_BET
	MTYPE_SITOUT_NEXT_ROUND
	MTYPE_JOIN_NEXT_ROUND
	// Mod ideas (not supported)
	MTYPE_SWAP_CARDS
	MTYPE_SWAP_SEATS
	MTYPE_BOMB_CARD
	// Other ideas: taxes
	// Other ideas: special combos
)

func MoveTypeStr2Int(mtype string) (int, error) {
	switch mtype {
	case "MTYPE_CHECK":
		return MTYPE_CHECK, nil
	case "MTYPE_FOLD":
		return MTYPE_FOLD, nil
	case "MTYPE_CALL":
		return MTYPE_CALL, nil
	case "MTYPE_CALL_ANY":
		return MTYPE_CALL_ANY, nil
	case "MTYPE_BET":
		return MTYPE_BET, nil
	case "MTYPE_SITOUT_NEXT_ROUND":
		return MTYPE_SITOUT_NEXT_ROUND, nil
	case "MTYPE_JOIN_NEXT_ROUND":
		return MTYPE_JOIN_NEXT_ROUND, nil
	case "MTYPE_SWAP_CARDS", "MTYPE_SWAP_SEATS", "MTYPE_BOMB_CARD":
		return 0, fmt.Errorf("Only vanilla poker supported")
	default:
		return 0, fmt.Errorf("Uknown move type")
	}
}

func MoveTypeInt2Str(mtype int) (string, error) {
	switch mtype {
	case MTYPE_CHECK:
		return "MTYPE_CHECK", nil
	case MTYPE_FOLD:
		return "MTYPE_FOLD", nil
	case MTYPE_CALL:
		return "MTYPE_CALL", nil
	case MTYPE_CALL_ANY:
		return "MTYPE_CALL_ANY", nil
	case MTYPE_BET:
		return "MTYPE_BET", nil
	case MTYPE_SITOUT_NEXT_ROUND:
		return "MTYPE_SITOUT_NEXT_ROUND", nil
	case MTYPE_JOIN_NEXT_ROUND:
		return "MTYPE_JOIN_NEXT_ROUND", nil
	case MTYPE_SWAP_CARDS, MTYPE_SWAP_SEATS, MTYPE_BOMB_CARD:
		return "", fmt.Errorf("Only vanilla poker supported")
	default:
		return "", fmt.Errorf("Uknown move type")
	}
}

// NOTE: GiveChips and Stop should be on queus to execute once certain conditions have passed
// NOTE: Only one game mode is currently supported: regular no-limits texas hold'em without mods or rake etc...
// NOTE: there is no thread safety; here is a list of locks we want
// 1. Lock for joining; joining should be atomic relative to private/public/changing bb and bb mult
// 2. Lock for moving; it should be atomic with respect to kicking, etc...
// 3. ...

// Defaults for standard games
const (
	DEFAULT_BB uint64                     = 1000               // Default big blind is 1000 chips
	DEFAULT_MAX_PLAYERS int               = 6                  // By default allow six players maximum
	DEFAULT_BB_MULTIPLIER uint64          = 100                // DEFAULT_BB * ..._MULTIPLIER = default starting hand
)

type Cards CardSet

func (c Cards) String() string {
	return CardSetToString(CardSet(c))
}

type Player struct {
	Name                string
	Hand                Cards
	Chips               uint64
	Bet                 uint64
	Pot                 uint64
	Admin               bool
	// Playing means keep them in the next round
	// PlayingRound means they are not yet folded
	Playing             bool
	PlayingRound        bool
	// SittingOutNextRound takes precedence over Playing
	SittingOutNextRound bool
	// Betting means they are at this moment poised to make a move (we are waiting on them)
	Betting             bool
	// exists is used to amortize player removal from the list
	exists              bool
}

type Pot struct {
	chips      uint64
	bettable   bool
	playerIdxs []int
}

type Game struct  {
	gameMaster    string
	streamCode    uuid.UUID
	name          string
	joinCode      string
	maxPlayers    int
	// Middle draw order is remembered by clients (remember game update messages)
	flop          Cards
	turn          Cards
	river         Cards
	middle        Cards
	// idx is the order in which they play
	players2idx   map[string]int
	players       []*Player
	// once only half the number of playersExist, players is modified to remove non-existing entries
	playersPlay   int
	playersExist  int
	private       bool
	// Paused means it's "frozen" and can recommence immediately on a request, while !started means
	// the game is not going on (i.e. it hasn't started yet, or it was stopped)
	paused        bool
	started       bool
	// The last pot(s) should be the betting pots
	pots          []*Pot
	bettingRound  int
	bettingIndex  int
	roundNum      int
	resolved      bool
	bigBlind      uint64
	// What to multiply by the bigBlind to get the number of startingChips
	bbMult     uint64
	// When someone all ins, if the next person bets over the newPotLimit,
	// a new pot is created
	someoneAllIn  bool
	newPotLimit   uint64
	// We send updates to the users via this channel
	gameUpdateCh  chan *GameUpdate
}

type GameUpdate struct {
	Message string
	// TODO (we are going to want to give more meaningful updates soon)
}

type GameInitArgs struct {
	// nil on any of these pointer values means "defaults please"
	Name          *string
	JoinCode      *string
	// Game Settings
	Private       *bool 
	MaxPlayers    *int
	BigBlind      *uint64
	BbMult        *uint64
}

// System maintenance; for now to change name or joinCode you'll need to make a new game
func New(creator string, args GameInitArgs) (string, *Game, error) {
	rand.Seed(time.Now().UnixNano())

	// We allow creators to request random name generation
	if creator == "" {
		creator = utils.RandPlayerName()
	}

	// Initialize game for streaming expecting to add a single admin player later
	g := &Game{
		gameMaster:   uuid.NewString(),
		streamCode:   uuid.New(),

		flop:         NoCards,
		turn:         NoCards,
		river:        NoCards,
		middle:       NoCards,

		players2idx:  make(map[string]int, 1),
		players:      make([]*Player, 1, 4),
		playersPlay:  1,
		playersExist: 1,
		started:      false,
		paused:       false,
		pots:         make([]*Pot, 0, 1),
		bettingRound: BROUND_NONE,
		bettingIndex: 0,
		roundNum:     0,

		resolved:     false,
		someoneAllIn: false,
		newPotLimit:  0,

		gameUpdateCh: make(chan *GameUpdate, 8),

	}

	// Initialize the GameInitArgs argumnts
	if args.Name != nil {
		g.name = *args.Name
	} else {
		g.name = fmt.Sprintf("%s-game-%s", creator, utils.RandString(3))
	}
	if args.JoinCode != nil {
		g.joinCode = *args.JoinCode
	} else {
		g.joinCode = fmt.Sprintf("%s-%s", utils.RandVerbAdv(nil), utils.RandString(3))
	}
	if args.Private != nil {
		g.private = *args.Private
	} else {
		g.private = false
	}
	if args.MaxPlayers != nil {
		g.maxPlayers = *args.MaxPlayers
	} else {
		g.maxPlayers = DEFAULT_MAX_PLAYERS
	}
	// Check the max players to make sure it's GUIable
	if g.maxPlayers < 0 {
		return "", nil, fmt.Errorf("Tried to set negative max players: %d", g.maxPlayers)
	}
	if g.maxPlayers > foreverMaxPlayers {
		return "", nil, fmt.Errorf("Tried to set max players above GUI limit: %d", g.maxPlayers)
	}
	if args.BigBlind != nil {
		g.bigBlind = *args.BigBlind
	} else {
		g.bigBlind = DEFAULT_BB
	}
	if args.BbMult != nil {
		g.bbMult = *args.BbMult
	} else {
		g.bbMult = DEFAULT_BB_MULTIPLIER
	}
	// Check the multipliers
	if maxUint64 / g.bigBlind < g.bbMult {
		return "", nil, fmt.Errorf("Cannot have bb %d and mult %d; overflows", g.bigBlind, g.bbMult)
	}
	
	// Add that single admin player (the creator)
	g.players[0] = &Player{
		Name:                creator,
		Hand:                NoCards,
		Chips:               g.bbMult * g.bigBlind,
		Bet:                 0,
		Pot:                 0,
		Admin:               true,
		Playing:             false,
		PlayingRound:        false,
		SittingOutNextRound: false,
		exists:              true,
	}
	g.players2idx[creator] = 0
	return creator, g, nil
}

// Helper functions
func (g *Game) findPlayer(p string) (*Player, bool) {
	idx, ok := g.players2idx[p]
	if !ok {
		return nil, false
	}
	pl := g.players[idx]
	if !pl.exists {
		return nil, false
	}
	return pl, true
}

func (g *Game) isAdmin(p string) (bool, error) {
	if p == g.gameMaster {
		return true, nil
	}
	pl, found := g.findPlayer(p)
	if !found {
		return false, fmt.Errorf("Could not find player %s", p)
	}
	return pl.Admin, nil
}

func (g *Game) ifAdmin(admin string, f func() (bool, error)) (bool, error) {
	mod, err := g.isAdmin(admin)
	if err != nil {
		return false, fmt.Errorf("Failed to check if player %s was admin: %v", admin, err)
	}
	if !mod {
		return false, nil
	}
	return f()
}

// Remove a player from the dictionary and players, keeping the order and doing the amortized
// delete if necesssary from the players list
func (g *Game) removePlayer(p string) error {
	idx, _ := g.players2idx[p]
	delete(g.players2idx, p)
	g.players[idx].exists = false
	g.playersExist -= 1
	if g.playersExist <= (len(g.players) >> 1) {
		resetBettingIdx := false
		
		newPlayers := make([]*Player, 0, g.playersExist)
		for idx, pl := range g.players {
			if pl.exists {
				newPlayers = append(newPlayers, pl)
			}
			if !resetBettingIdx && idx == g.bettingIndex {
				if !pl.exists {
					return fmt.Errorf("Somehow betting index as on non-existent player")
				}
				g.bettingIndex = len(newPlayers) - 1
				resetBettingIdx = true
			}
			g.players2idx[pl.Name] = len(newPlayers) - 1
		}
		g.players = newPlayers
	}
	return nil
}

// A player can only be kicked if they're not betting/making a move
// in the future we'll use a mutex to make sure this is atomic
// and thus we can kick while betting without risking them making a move at the same time
func (g *Game) KickPlayer(kicker string, kicked string) (bool, error) {
	return g.ifAdmin(kicker, func() (bool, error) {
		pl, found := g.findPlayer(kicked)
		if !found {
			return false, fmt.Errorf("Player %s was not found", kicked)
		}
		if pl.Betting {
			return false, fmt.Errorf("Cannot kick %s while he's betting", kicked)
		}
		return true, g.removePlayer(kicked)
	})
}

// A player can be modded or unmodded at any time
func (g *Game) ModPlayer(modder string, modded string, admin bool)  (bool, error) {
	return g.ifAdmin(modder, func() (bool, error) {
		pl, found := g.findPlayer(modded)
		if !found {
			return false, fmt.Errorf("Tried to mod nonexistent player: %s", modded)
		}
		pl.Admin = admin
		return true, nil
	})
}

// A player can be added at any moment in time, but they won't play until the next round
func (g *Game) AddPlayer(newPlayer string, joinCode string) (string, bool, error) {
	if g.private && joinCode != g.joinCode {
		return "", false, nil
	}
	_, found := g.findPlayer(newPlayer)
	if found {
		return "", false, fmt.Errorf("Player with name %s was already found in game", newPlayer)
	}
	if newPlayer == "" {
		newPlayer = utils.RandPlayerName()
	}
	pl := &Player{
		Name:                newPlayer,
		Hand:                NoCards,
		Chips:               g.bigBlind * g.bbMult,
		Bet:                 0,
		Pot:                 0,
		Admin:               false,
		Playing:             true,
		PlayingRound:        false,
		SittingOutNextRound: false,
		Betting:             false,
		exists:              true,
	}
	g.players2idx[newPlayer] = len(g.players)
	g.players = append(g.players, pl)
	g.playersExist += 1
	// In the future we may want to pick a random location between
	// another pair of players and insert there; however, it would take
	// linear time unless we could guarantee no more than half of the players
	// existed (or something like that) (note: you can't just swap because that
	// changes the rest of the order)

	return newPlayer, true, nil
}

// We return a copy so that people can't modify this dangerously
// We delete the players who don't exist
func (g *Game) Players() []*Player {
	players := make([]*Player, 0, len(g.players))
	for _, pl := range g.players {
		if pl.exists {
			players = append(players, pl)
		}
	}
	return players
}

func (g *Game) BigBlind() uint64 {
	return g.bigBlind
}

// Remember that Cards is just a uint64 and thus has no order, but is safe to return
func (g *Game) Middle() Cards {
	return g.middle
}

// We make a copy to make sure that this can't be modified dangerously
func (g *Game) Pots() []uint64 {
	p := make([]uint64, len(g.pots))
	for i, pot := range g.pots {
		p[i] = pot.chips
	}
	return p
}

// Game Status

// Play and pause are instantaneous so they are ok to go through now; their only
// purpose is to unfreeze/freeze any timers and unlock/lock the game state blocking any moves until
// the next play
func (g *Game) Play(requester string) (bool, error) {
	return g.ifAdmin(requester, func() (bool, error) {
		g.paused = false
		return true, nil
	})
}

func (g *Game) Pause(requester string) (bool, error) {
	return g.ifAdmin(requester, func() (bool, error) {
		g.paused = true
		return true, nil
	})
}

func (g *Game) Start(requester string) (bool, error) {
	return g.ifAdmin(requester, func() (bool, error) {
		// Need at least two players to play this game
		if g.playersPlay < 2 {
			return false, nil
		}
		if g.started {
			return false, fmt.Errorf("Game already started")
		}
		g.started = true
		// Incrementing from BROUND_NONE hands out cards, etc...
		return g.Increment()
	})
}

func (g *Game) Stop(requester string) (bool, error) {
	return g.ifAdmin(requester, func() (bool, error) {
		if !g.started {
			return false, fmt.Errorf("Game was never started")
		}
		// In reality we will probably want to have this be a queue for players
		// we may also want to allow stopping mid-round if paused or something
		// else like that
		if g.bettingRound != BROUND_NONE {
			return false, nil
		}
		g.started = false
		// Clean up... luckily Increment() and NewRound() should do much of the heavy lifting for us 
		// once BROUND_NONE is reached since at that point we are just waiting for another Increment()
		// to start a new round (with all its inner betting rounds).
		g.paused = false
		return true, nil
	})
}


// This is ok to do at any time since it only affects people who are joining
func (g *Game) MakePrivate(requester string) (bool, error) {
	return g.ifAdmin(requester, func() (bool, error) {
		g.private = true
		return true, nil
	})
}

// This is ok to do at any time since it only affects people who are joining
func (g *Game) MakePublic(requester string) (bool, error) {
	return g.ifAdmin(requester, func() (bool, error) {
		g.private = false
		return true, nil
	})
}

// This is ok to do at any time since it only affects people who are joining
func (g *Game) ChangeMaxPlayers(requester string, newMaxPlayers int) (bool, error) {
	return g.ifAdmin(requester, func() (bool, error) {
		if newMaxPlayers < 0 {
			return false, fmt.Errorf("Cannot use negative max num players: %d", newMaxPlayers)
		}
		if newMaxPlayers > foreverMaxPlayers {
			return false, fmt.Errorf("max players must be under %d which is true for ALL games", foreverMaxPlayers)
		}
		g.maxPlayers = newMaxPlayers
		return true, nil
	})
}

// This can be done at any time since it only affects players who are joining
func (g *Game) ChangeBb(requester string, newBb uint64) (bool, error) {
	return g.ifAdmin(requester, func() (bool, error) {
		if maxUint64 / g.bbMult < newBb {
			return false, fmt.Errorf("New big blind %d would overflow bb mult %d", newBb, g.bbMult)
		}
		g.bigBlind = newBb
		return true, nil
	})
}

// This can be changed immediately because it only affects players who may join
func (g *Game) ChangeBbMult(requester string, newMult uint64) (bool, error) {
	return g.ifAdmin(requester, func() (bool, error) {
		if maxUint64 / g.bigBlind < newMult {
			return false, fmt.Errorf("New mult %d would overflow big blind %d", newMult, g.bigBlind)
		}
		g.bbMult = newMult
		return true, nil
	})
}

func (g *Game) Started() bool {
	return g.started
}

func (g *Game) Paused() bool {
	return g.paused
}

func (g *Game) Private() bool {
	return g.private
}

func (g *Game) Name() string {
	return g.name
}

func (g *Game) RoundsPlayed() int {
	return g.roundNum
}

func (g *Game) StreamCode() uuid.UUID {
	return g.streamCode
}

func (g *Game) NewStreamCode() uuid.UUID {
	g.streamCode = uuid.New()
	return g.StreamCode()
}

func (g *Game) GameUpdateChan() chan *GameUpdate {
	return g.gameUpdateCh
}

// Game Flow
func (g *Game) Move(move int, chips uint64, mover string) (bool, error) {
	// TODO remember to turn of pl.betting once they bet
	return false, fmt.Errorf("Not Implemented")
}

// A name can be changed at any time, since we use the idx in the Players list as an id
func (g *Game) ChangePlayerName(changer string, changed string, newName string) (string, bool, error) {
	if newName == "" {
		newName = utils.RandPlayerName()
	}
	ch, err := g.ifAdmin(changer, func() (bool, error) {
		pl, found := g.findPlayer(changed)
		if !found {
			return false, fmt.Errorf("Tried to change player %s name, but wasn't found", changed)
		}
		// Since pl is a pointer to a player, this modifies the underlying g.players list
		pl.Name = newName
		g.players2idx[newName] = g.players2idx[changed]
		delete(g.players2idx, changed)
		return true, nil
	})
	return newName, ch, err
}

// Chips are OK to give at any time, so long as the player is not playing this round
// (later we may want to add this to a queue instead)
func (g *Game) GiveChips(giver string, reciever string, amount uint64) (bool, error) {
	return g.ifAdmin(giver, func() (bool, error) {
		pl, found := g.findPlayer(reciever)
		if !found {
			return false, fmt.Errorf("Failed to find player %s to give chips to", reciever)
		}
		if maxUint64 - amount > pl.Chips {
			return false, fmt.Errorf("Too many chips: started with %d and tried to add %d, but that would overflow (max: %d)", pl.Chips, amount, maxUint64)
		}
		if pl.PlayingRound {
			return false, nil
		}
		pl.Chips += amount
		return true, nil
	})
}

// Game Flow Control Plane

// Helpers...
func (g *Game) incrementBettingIndex() {
	g.bettingIndex ++
	for !g.players[g.bettingIndex].exists && g.bettingIndex < len(g.players) {
		g.bettingIndex ++
	}
}

// hand out all the cards to middle and players too
// it could use certain efficiency improvements
func (g *Game) handout() {
	cards := make([]Cards, 52)
	cards[0] = Cards(AceOfClubs)
	for i := 1; i < 4; i++ {
		cards[i] = cards[i-1] << 1
	}
	cards[4] = Cards(TwoOfClubs)
	for i := 5; i < 52; i++ {
		cards[i] = cards[i-1] << 1
	}
	rand.Shuffle(2 * len(cards), func(i int, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})

	i := 0
	for _, pl := range g.players{
		if pl.Playing {
			pl.Hand = cards[i] | cards[i + 1]
			i ++
		}
	}
	g.flop = cards[i] | cards[i + 1] | cards[i + 2]
	g.turn = cards[i + 3]
	g.river = cards[i + 4]
}

// TODO increment betting index conditionally

// We use bettingIndex = len(g.players) as a convention that says "ready to move on to the next
// betting round". Increment is meant to be called on each state change, so any time
// some callback or incoming request forces us to change state, we do so by using Increment()
func (g *Game) Increment() (bool, error) {
	switch g.bettingRound {
	case BROUND_NONE:
		// Give everyone their cards
		g.handout()
		g.bettingRound = BROUND_PREFLOP
	case BROUND_PREFLOP:
		if g.bettingIndex == len(g.players) {
			// TODO
		}
	case BROUND_FLOP:
		if g.bettingIndex == len(g.players) {
			// TODO
		}
		// TODO
	case BROUND_TURN:
		if g.bettingIndex == len(g.players) {
			// TODO
		}
		// TODO
	case BROUND_RIVER:
		if g.bettingIndex == len(g.players) {
			if g.resolved {
				return g.NewRound()
			}
			return g.Resolve()
		}
		// Once you place a bet you are done betting
		if !g.players[g.bettingIndex].Betting {
			g.incrementBettingIndex()
			// TODO
		}
		// TODO need to clean up bettable pots
		// TODO need to clean up someone all in and newPotLimit
		// TODO need to clean up bets
	default:
		return false, fmt.Errorf("Unknown betting round")
	}
	return false, nil
}

// Return the highest 
func highestCardset(s CardSet) (CardSet, int) {
	if c := royalFlush(s); c > 0 {
		return c, HTYPE_ROYAL_FLUSH
	}
	if c := fourOfAKind(s); c > 0 {
		return c, HTYPE_QUAD
	}
	if c := straightFlush(s); c > 0 {
		return c, HTYPE_STRAIGHT_FLUSH
	}
	if c := fullHouse(s); c > 0 {
		return c, HTYPE_FULLHOUSE
	}
	if c := flush(s); c > 0 {
		return c, HTYPE_FLUSH
	}
	if c := straight(s); c > 0 {
		return c, HTYPE_STRAIGHT
	}
	if tHigh, tLow := triplet(s); tHigh > 0 {
		return tHigh | tLow, HTYPE_TRIP
	}
	if pHigh, pMed, pLow := pair(s); pHigh > 0 {
		if pMed > 0 {
			return pHigh | pMed | pLow, HTYPE_TWO_PAIR
		}
		return pHigh, HTYPE_PAIR
	}
	return s, HTYPE_HIGHCARD
}

func (g *Game) winner(a *Player, b *Player) (*Player, Cards, int) {
	csA, htypeA := highestCardset(CardSet(a.Hand | g.middle))
	csB, htypeB := highestCardset(CardSet(b.Hand | g.middle))

	// NOTE we are allowing winners to win by suit, which may not always be
	// desireable (in most online sites this is not the case)
	if htypeB > htypeA || (htypeB == htypeA && csB > csA) {
		return b, Cards(csB), htypeB
	} else if htypeB == htypeA && csA == csB {
		return nil, Cards(csA), htypeA
	}
	return a, Cards(csA), htypeA
}

// Needs to deal with multiple pots (etc...)
func (g *Game) Resolve() (bool, error) {
	if g.resolved {
		return false, fmt.Errorf("Cannot call resolve when the game is already resolved")
	}
	// For each pot give the pot to the winner
	for _, pot := range g.pots {
		winners := make([]*Player, 1)
		if pot.bettable {
			return false, fmt.Errorf("Pot was bettable when it should not have been")
		}
		if len(pot.playerIdxs) < 1 {
			return false, fmt.Errorf("Empty player idxs in this pot")
		}
		var pl *Player
		for _, idx := range pot.playerIdxs {
			if pl == nil {
				pl = g.players[idx]
				winners[0] = pl
			} else {
				// NOTE: we may want to send the hand type on the channel
				w, _, _ := g.winner(pl, g.players[idx])
				if w == nil {
					winners = append(winners, g.players[idx])
				}
			}
		}
		if pl == nil {
			return false, fmt.Errorf("Somehow pl was nil")
		}
		// NOTE: there is a rounding mistake here, though it is very minor
		for _, w := range winners {
			w.Chips += (pot.chips / uint64(len(winners))) + 1
		}
	}

	// Clean up player financials
	for _, pl := range g.players {
		pl.Pot = 0
	}

	g.resolved = true
	return true, nil
}

func (g *Game) NewRound() (bool, error) {
	if !g.resolved {
		return false, fmt.Errorf("Tried to start a new round without resolving the old one")
	}
	// Cleanup round
	g.bettingRound = BROUND_NONE

	g.middle = NoCards
	g.flop = NoCards
	g.turn = NoCards
	g.river = NoCards
	g.pots = make([]*Pot, 0, 1)

	g.bettingIndex = 0
	g.incrementBettingIndex()

	g.roundNum += 1

	// Cleanup players
	for _, pl := range g.players {
		pl.Hand = NoCards
		if pl.SittingOutNextRound {
			pl.SittingOutNextRound = false
			pl.Playing = false
			pl.PlayingRound = false
			g.playersPlay -= 1
		}
	}

	if g.playersPlay < 2 {
		st, err := g.Stop(g.gameMaster)
		if !st || err != nil {
			return false, fmt.Errorf("Failed to stop game at under 2 players: stop = %v, err = %v", st, err)
		}
	}

	g.resolved = false
	return true, nil
}