package poker

import "testing"

import (
	_ "github.com/4gatepylon/GoPoker/poker"
)

// Golang is smart and every file that is Test<string that doesn't start with a lowercase letter>
// will be run when you do "go test" and we can use fmt.Errof to throw errors

// I might upgrade to Bazel later (though I can write my own rule in that case....)

func TestRoyalFlushExist(t *testing.T) {
	// try all the royal flushes
	// try more than one royal flush at a time (make sure highest returns)
	// try additional cards (1 and all) and make sure that it works regardless
	var royalFlushesInput = [7]CardSet{
		SpadesRoyalFlush,
		HeartsRoyalFlush,
		DiamondsRoyalFlush,
		ClubsRoyalFlush,
		HeartsRoyalFlush | FiveOfClubs,
		DiamondsRoyalFlush | TenOfSpades,
		AllCards,
	}

	var royalFlushesExpect = [7]CardSet{
		SpadesRoyalFlush,
		HeartsRoyalFlush,
		DiamondsRoyalFlush,
		ClubsRoyalFlush,
		HeartsRoyalFlush,
		DiamondsRoyalFlush,
		SpadesRoyalFlush,
	}

	for i := 0; i < 7; i++ {
		var royalFlushExpect CardSet = royalFlushesExpect[i]
		var royalFlushResult CardSet = royalFlush(royalFlushesInput[i])

		if royalFlushExpect != royalFlushResult {
			t.Errorf(
				"\nGOT\n{%s},\nbut EXPECTED\n{%s}", 
				CardSetToString(royalFlushResult), 
				CardSetToString(royalFlushExpect),
			)
		}
	}
}

// func TestRoyalFlushDoesNotExist(t *testing.T) {
// 	t.Errorf("cunt")
// }

// func TestQuadsExist(t *testing.T) {
// 	t.Errorf("cunt")
// }

// func TestQuadsDoNotExist(t *testing.T) {
// 	t.Errorf("cunt")
// }

func TestCardSetToString(t *testing.T) {
	t.Errorf("NOT IMPLEMENTED!\nAll cards below.\n%s", CardSetToString(AllCards))
}