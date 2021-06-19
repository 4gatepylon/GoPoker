package poker

import "testing"
import "fmt"
import "strings"

// UTILITY
func gotButExpected(got string, expected string) string {
	return fmt.Sprintf("\nGOT\n{%s},\nbut EXPECTED\n{%s}", got, expected)
}

// short name purely for practical reasons
func errMsg(got CardSet, expected CardSet) string {
	return gotButExpected(CardSetToString(got), CardSetToString(expected))
}

// so we can have quick contains checks basically
type StringSet map[string]struct{}

func strSet(elements... string) StringSet {
	st := make(StringSet, 2)
	for _, k := range elements {
		st[k] = struct{}{}
	}
	return st
}

func (s StringSet) contains(e string) bool {
	var _, contained = s[e]
	return contained
}

func strSet2str(s StringSet) string {
	var b strings.Builder
	b.WriteString("{ ")
	for k, _ := range s {
		b.WriteString(fmt.Sprintf("\"%s\", ", k))
	}
	b.WriteString("}")

	return b.String()
}

// TESTS
func TestRoyalFlushExists(t *testing.T) {
	const numTests = 7

	// try all the royal flushes
	// try more than one royal flush at a time (make sure highest returns)
	// try additional cards (1 and all) and make sure that it works regardless
	var royalFlushesInput = [numTests]CardSet{
		SpadesRoyalFlush,
		HeartsRoyalFlush,
		DiamondsRoyalFlush,
		ClubsRoyalFlush,
		HeartsRoyalFlush | FiveOfClubs,
		DiamondsRoyalFlush | TenOfSpades,
		AllCards,
	}

	var royalFlushesExpect = [numTests]CardSet{
		SpadesRoyalFlush,
		HeartsRoyalFlush,
		DiamondsRoyalFlush,
		ClubsRoyalFlush,
		HeartsRoyalFlush,
		DiamondsRoyalFlush,
		SpadesRoyalFlush,
	}

	for i := 0; i < numTests; i++ {
		var royalFlushExpect CardSet = royalFlushesExpect[i]
		var royalFlushResult CardSet = royalFlush(royalFlushesInput[i])

		if royalFlushExpect != royalFlushResult {
			t.Errorf(errMsg(royalFlushResult, royalFlushExpect))
		}
	}
}

func TestRoyalFlushDoesNotExist(t *testing.T) {
	const numTests = 5

	var notRoyalFlushesInput = [numTests]CardSet{
		// test with only one missing
		SpadesRoyalFlush & ^KingOfSpades,

		//...
		NoCards,

		// straight no flush
		AceOfDiamonds | KingOfDiamonds | QueenOfClubs | JackOfClubs | TenOfHearts, 
		
		// flush no straight
		Clubs & ^AceOfClubs, 

		//arbitrary cards
		ThreeOfHearts | TenOfClubs | SevenOfDiamonds | EightOfSpades | AceOfSpades | AceOfDiamonds | AceOfHearts, 
	}

	var expectedOutput CardSet = 0

	for i := 0; i < numTests; i++ {
		var output CardSet = royalFlush(notRoyalFlushesInput[i])

		if output != expectedOutput {
			t.Errorf(errMsg(output, expectedOutput))
		}
	}
}

// func TestQuadsExist(t *testing.T) {
// 	t.Errorf("Unimplemented!")
// }

// func TestQuadsDoNotExist(t *testing.T) {
// 	t.Errorf("Unimplemented!")
// }

// func TestStraightExists(t *testing.T) {
// 	t.Errorf("Unimplemented!")
// }

// func TestStaightDoesNotExist(t *testing.T) {
// 	t.Errorf("Unimplemented!")
// }

// func TestStraightFlushExists(t *testing.T) {
// 	t.Errorf("Unimplemented!")
// }

// func TestStaightFlushDoesNotExist(t *testing.T) {
// 	t.Errorf("Unimplemented!")
// }

// func TestFlushExists(t *testing.T) {
// 	t.Errorf("Unimplemented!")
// }

// func TestFlushDoesNotExist(t *testing.T) {
// 	t.Errorf("Unimplemented!")
// }

// func TestFullHouseExists(t *testing.T) {
// 	t.Errorf("Unimplemented!")
// }

// func TestFullHouseDoesNotExist(t *testing.T) {
// 	t.Errorf("Unimplemented!")
// }

// func TestTripsExist(t *testing.T) {
// 	t.Errorf("Unimplemented!")
// }

// func TestTripsDoNotExist(t *testing.T) {
// 	t.Errorf("Unimplemented!")
// }

// func TestPairsExists(t *testing.T) {
// 	t.Errorf("Unimplemented!")
// }

// func TestPairsDoNotExist(t *testing.T) {
// 	t.Errorf("Unimplemented!")
// }

// func TestHighCard(t *testing.T) {
// 	// test with no cards
// 	// test with a high card
// 	t.Errorf("Unimplemented!")
// }


// tests that the proper strings are displayed
// the proper seperator is the space, cards are expected to be sorted
// by numerical order (increasing) but not by suit order
func TestCardSetToString(t *testing.T) {
	const numTests = 8

	var cardsets = [numTests]CardSet{
		// a couple singletons
		AceOfSpades,
		JackOfDiamonds,
		ThreeOfClubs,
		// a couple non-increasing sets
		TenOfClubs | TenOfDiamonds,
		SevenOfHearts | SevenOfSpades | SevenOfDiamonds,
		// a couple increasing sets
		FourOfDiamonds | FiveOfClubs,
		// increasing and non-increasing sets (i.e. subsets exist that increase and that don't increase)
		EightOfSpades | NineOfClubs | NineOfHearts,
		// empty should be empty string
		NoCards,
	}

	var validAnswers = [numTests]StringSet {
		strSet("AS"),
		strSet("JD"),
		strSet("3C"),
		strSet("TC TD", "TD TC"),
		strSet("7H 7S 7D", "7H 7D 7S", "7D 7S 7H", "7D 7H 7S", "7S 7H 7D", "7S 7D 7H"),
		strSet("4D 5C"),
		strSet("8S 9C 9H", "8S 9H 9C"),
		strSet(""),
	}

	for i := 0; i < numTests; i++ {
		var output string = CardSetToString(cardsets[i])

		if !validAnswers[i].contains(output) {
			t.Errorf(gotButExpected(output, fmt.Sprintf("One of %s", strSet2str(validAnswers[i]))))
		}
	}
}