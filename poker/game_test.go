package poker

import (
	"testing"
)

// The same creator name is used for all tests for convenience
const creator string = "creator"

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

func TestKickAllowed(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}

func TestKickNotAllowed(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}

func TestModAllowed(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}

func TestModNotAllowed(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}

func TestAddPlayerPublic(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}

func TestAddPlayerPrivateJoinCodeOk(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}

func TestAddPlayerPrivateJoinCodeNotOk(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}

func TestPlayersPlayingMany(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}

func TestPlayersNotPlayingOne(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
	}, t)
}

func TestPlayersPausedMany(t *testing.T) {
	wrap(func(game GameLike, t *testing.T) {
		// TODO test
	}, New, creator, &GameInitArgs{
		// TODO args
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

func TestPlayPauseAllowed(t *testing.T) {
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
