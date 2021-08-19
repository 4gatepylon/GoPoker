package poker

import (
	"testing"
	"fmt"
)

// TODO finish this unit-testing (at least in a single-threaded environment)

// The same creator name is used for all tests for convenience
const creator string = "creator"
const join_code string = "join_code"
const false_join_code string = "false_join_code"
const game_name string = "game"

func pointer(c string) *string {
	return &c
}

// wrap() handles setup and teardown for tests.
// The idea of this framework is to be implementation
// agnostic. You declare a TestX(t *testing.T) function
// as well as a func(game GameLike, *testing.T) function
// which tests the GameLike interface implementation of your choice.
// You choose my passing in an init function which returns the
// type which implements the interface. Additionally, note that Teardown()
// will be tested per-implementation. These tests assume Teardown() works.

func wrap(
	test func(GameLike, *testing.T),
	init func(*string, *GameInitArgs) (GameLike, error),
	creator string,
	args *GameInitArgs,
	t *testing.T) {
	game, err := init(&creator, args)
	if err != nil {
		t.Fatalf("Error initializing game: `%v`\n", err)
	}
	defer game.Teardown()
	test(game, t)
}

func TestAddPlayersPublicKickAllowedAndNotAllowed(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// Test able to add
		noob := "noob"
		pres, added, err := game.AddPlayer(&noob, nil)
		if !added || err != nil || pres == nil {
			t.Fatalf("Failed to add (added: %v) player (resulting name pointer %d) with err: `%v`", added, pres, err)
		}
		if *pres != noob {
			t.Fatalf("Tried to add `%v` but added `%v`", noob, *pres)
		}

		// Test able to kick
		kicker := "kicker"
		_creator := creator
		game.AddPlayer(&kicker, nil)

		kicked, _ := game.KickPlayer(&kicker, &noob)
		if kicked {
			t.Fatalf("Kicked when should not be allowed")
		}
		kicked, err = game.KickPlayer(&_creator, &noob)
		if err != nil || !kicked {
			t.Fatalf("Failed to kick (kicked: %v) when should be able to (err: `%v`)", kicked, err)
		}
	}, New, creator, &GameInitArgs{
		Name: pointer(game_name),
		Public: true,
	}, t)
}

func TestAddPlayerPrivateJoinCodeOkAndNotOk(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		friend := "friend"
		foe := "foe"
		_, joined, err := game.AddPlayer(&friend, pointer(join_code))
		if !joined || err != nil {
			t.Fatalf("Failed to add player (joined: %v) when should have been possible, err: `%v`", joined, err)
		}
		_, joined, err = game.AddPlayer(&foe, pointer(false_join_code))
		if joined || err != nil {
			t.Fatalf("Added player (joined: %v) when shouldn't have (err: `%v`)", joined, err)
		}
	}, New, creator, &GameInitArgs{
		Name: pointer(game_name),
		Public: false,
		JoinCode: pointer(join_code),
	}, t)
}

func TestModAllowedAndNotAllowed(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		friend := "friend"
		foe := "foe"
		minion := "minion"
		game.AddPlayer(&friend, nil)
		game.AddPlayer(&foe, nil)
		game.AddPlayer(&minion, nil)

		modded, err := game.ModPlayer(pointer(creator), &friend, PPERM_ADMIN)
		if !modded || err != nil {
			t.Fatalf("Modded %v, while err was `%v`, but expected to mod without an error", modded, err)
		}
		modded, err = game.ModPlayer(pointer(foe), &minion, PPERM_ADMIN)
		if modded {
			t.Fatalf("Modded when should have not: `%v`", err)
		}
	}, New, creator, &GameInitArgs{
		Name: pointer(game_name),
		Public: true,
	}, t)
}

func TestPublicPlayersPlayingStakesManyPlayPauseAllowed(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// Initialize with creator + p1,p2,p3,...,p5
		found := map[string]bool{creator: false}
		for i := 1; i <= 5; i++ {
			p := fmt.Sprintf("p%d", i)
			game.AddPlayer(pointer(p), nil)
			found[p] = false
		}

		// Make sure you can play
		playing, err := game.Play(pointer(creator))
		if playing_check := game.Playing(); !playing || !playing_check || err != nil {
			t.Fatalf("Failed to start playing (playing first %v, then %v): `%v`", playing, playing_check, err)
		}
		players := game.Players()
		if players == nil || len(players) < 6 {
			t.Fatalf("Got %d players but should have been 6", len(players))
		}

		// Make sure each is of the expected ones
		for i := 0; i < len(players); i++ {
			p := players[i].Name
			previously, ok := found[p]
			if previously || ok {
				t.Fatalf("Previously: %v, ok: %v, shold have been not ok", previously, ok)
			}
			found[p] = true
		}

		// Make sure the stakes are ok
		if s := game.Stakes(); s != 1000 {
			t.Fatalf("Game has stakes %d but should be %d", s, 1000)
		}

		// Make sure you can pause
		paused, err := game.Pause(pointer(creator))
		if paused_check := game.Playing(); !paused || paused_check || err != nil {
			t.Fatalf("Failed to pause (first with %v then with %v) with err `%v`", paused, paused_check, err)
		}
	}, New, creator, &GameInitArgs{
		Name: pointer(game_name),
		Public: true,
		Stakes: 1000,
	}, t)
}

func TestPlayersNotPlayingOne(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		Name: pointer(game_name),
		Public: true,
	}, t)
}

func TestPlayersPausedMany(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		Name: pointer(game_name),
		Public: true,
	}, t)
}

func TestStakesPlayingAndChanging(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}

func TestStakesNotPlayingAndChanging(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}

func TestMiddleEmpty(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}

func TestMiddleFlopTurnAndRiver(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}

func TestPotsBeforePlay(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}

func TestPotsOneZeroOrNoneInBettingRound(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}

func TestPotsOneInBettingRound(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}

func TestPotsManyInBettingRound(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}

func TestChangeGameNameAllowed(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}

func TestChangeGameNameNotAllowed(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}

func TestPlayPauseNotAllowed(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}

func TestMakePublicPrivateAllowed(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}

func TestPublicPrivateNotAllowed(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}

func TestPlayingPrivate(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}

func TestMoveCheckIsTurn(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}

func TestMoveCallIsTurn(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}

func TestMoveBetIsTurnAndAllIn(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}

func TestMoveIsTurnPaused(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}

func TestMoveIsNotTurnPlayAndPause(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}

func TestMoveSitOutNextRound(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}

func TestMoveCallAnyAndAllIn(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}

func TestChangePlayerNameAllowedModAndSelf(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}

func TestChangePlayerNameNotAllowed(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}

func TestIncrementAllowed(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}

func TestIncrementNotAllowed(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}

func TestResolveOnePotOneWinner(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}

func TestResolveOnePotTwoWinnersTie(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}

func TestResolveManyPotsOneWinner(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}

func TestResolveThreeWayTieManyPots(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}

func TestNewRoundAllowed(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}

func TestNewRoundNotAllowed(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}

func TestRenew(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}

func TestRenewKeepChipsAndPlayers(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}

func TestRenewKeepPlayers(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}
