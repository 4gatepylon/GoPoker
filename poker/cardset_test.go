package poker

import (
	"testing"
	"fmt"
	"strings"
)

// yes I know there is a lot of code copying, sorry m88

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

func TestQuadsExist(t *testing.T) {
	const numTests = 5

	var quads = [numTests]CardSet{
		// check with some quads
		Aces, // check aces always as they are an odd case
		Twos,
		Tens,
		// check with more than one quad (ya never gonna happen but who knows you might mod the game so you're welcome)
		Threes | Fours,
		Sevens | Kings,
	}

	var expectedOutputs = [numTests]CardSet{
		Aces,
		Twos,
		Tens,
		Fours,
		Kings,
	}

	for i := 0; i < numTests; i++ {
		var expectedOutput = expectedOutputs[i]
		
		var output CardSet = fourOfAKind(quads[i])

		if output != expectedOutput {
			t.Errorf(errMsg(output, expectedOutput))
		}
	}
}

func TestQuadsDoNotExist(t *testing.T) {
	const numTests = 6

	var notQuads = [numTests]CardSet{
		// check a couple with one missing
		Aces & ^AceOfDiamonds,
		Nines & ^NineOfHearts,
		// a couple with two missing
		KingOfHearts | KingOfSpades | QueenOfSpades | TenOfHearts | SixOfHearts | SixOfDiamonds | TwoOfClubs,
		// a couple with randomness
		Clubs,
		// a couple with two 
		Queens & Jacks,
		NoCards,
	}

	var expectedOutput CardSet = NoCards

	for i := 0; i < numTests; i++ {
		var output CardSet = fourOfAKind(notQuads[i])
		if output != expectedOutput {
			t.Errorf(errMsg(output, expectedOutput))
		}
	}
}

func TestStraightExists(t *testing.T) {
	const numTests = 5

	var theres = [numTests]CardSet{
		// top down and down up with spades
		AceOfSpades | KingOfClubs | QueenOfDiamonds | JackOfHearts | TenOfHearts,
		FiveOfDiamonds | FourOfClubs | ThreeOfHearts | TwoOfHearts | AceOfHearts,
		// inefficient
		TenOfClubs | NineOfClubs | EightOfClubs | SevenOfClubs | SixOfClubs | FiveOfDiamonds,
		QueenOfHearts | JackOfDiamonds | JackOfSpades | TenOfDiamonds | NineOfDiamonds | EightOfClubs,
		AceOfSpades | TenOfSpades | NineOfSpades | SevenOfSpades | SixOfSpades | EightOfClubs,
	}

	var expectedOutputs = [numTests]CardSet{
		AceOfSpades | KingOfClubs | QueenOfDiamonds | JackOfHearts | TenOfHearts,
		FiveOfDiamonds | FourOfClubs | ThreeOfHearts | TwoOfHearts | AceOfHearts,
		TenOfClubs | NineOfClubs | EightOfClubs | SevenOfClubs | SixOfClubs,
		// NOTE: pairs WILL be returned, which should be OK, since we only care about the highest card
		QueenOfHearts | JackOfDiamonds | JackOfSpades | TenOfDiamonds | NineOfDiamonds | EightOfClubs,
		TenOfSpades | NineOfSpades | SevenOfSpades | SixOfSpades | EightOfClubs,
	}

	for i := 0; i < numTests; i++ {
		var expectedOutput CardSet = expectedOutputs[i]
		var output CardSet = straight(theres[i])

		if output != expectedOutput {
			fmt.Println(output, expectedOutput)
			t.Errorf(errMsg(output, expectedOutput))
		}
	}
}

func TestStaightDoesNotExist(t *testing.T) {
	const numTests = 5

	var notTheres = [numTests]CardSet{
		// off by one
		AceOfSpades | KingOfClubs | QueenOfSpades | JackOfDiamonds | NineOfDiamonds,
		FiveOfClubs | FourOfClubs | ThreeOfClubs | TwoOfDiamonds | TenOfHearts,
		FourOfClubs | AceOfClubs | TwoOfHearts | ThreeOfHearts,
		// totally off
		JackOfClubs | JackOfDiamonds | JackOfHearts | TenOfClubs | NineOfHearts | EightOfSpades,
		NoCards,
	}

	var expectedOutput CardSet = 0

	for i := 0; i < numTests; i++ {
		var output CardSet = straight(notTheres[i])

		if output != expectedOutput {
			t.Errorf(errMsg(output, expectedOutput))
		}
	}
}

func TestStraightFlushExists(t *testing.T) {
	const numTests = 3

	var theres = [numTests]CardSet{
		// excess
		Spades,
		// exact
		(Tens & ^ TenOfClubs) | NineOfDiamonds | EightOfDiamonds | SevenOfDiamonds | SixOfDiamonds,
		FiveOfHearts | FourOfHearts | ThreeOfHearts | TwoOfHearts | AceOfHearts,
	}

	var expectedOutputs = [numTests]CardSet{
		AceOfSpades | KingOfSpades | QueenOfSpades | JackOfSpades | TenOfSpades,
		TenOfDiamonds | NineOfDiamonds | EightOfDiamonds | SevenOfDiamonds | SixOfDiamonds,
		FiveOfHearts | FourOfHearts | ThreeOfHearts | TwoOfHearts | AceOfHearts,
	}

	for i := 0; i < numTests; i++ {
		var expectedOutput CardSet = expectedOutputs[i]
		var output CardSet = straightFlush(theres[i])

		if output != expectedOutput {
			t.Errorf(errMsg(output, expectedOutput))
		}
	}
}

func TestStaightFlushDoesNotExist(t *testing.T) {
	const numTests = 5

	var notTheres = [numTests]CardSet{
		// straight and flush but not straight flush
		AceOfSpades | KingOfSpades | QueenOfClubs | JackOfSpades | TenOfSpades | NineOfSpades,
		// straight but not flush
		SixOfClubs | FiveOfHearts | FourOfHearts | ThreeOfDiamonds | TwoOfSpades | AceOfSpades,
		FiveOfClubs | AceOfDiamonds | FourOfClubs | ThreeOfClubs | TwoOfClubs | JackOfDiamonds,
		// flush but not straight
		Diamonds & ^TenOfDiamonds & ^FiveOfDiamonds,
		// nothing
		NoCards,
	}

	var expectedOutput CardSet = 0

	for i := 0; i < numTests; i++ {
		var output CardSet = straightFlush(notTheres[i])

		if output != expectedOutput {
			t.Errorf(errMsg(output, expectedOutput))
		}
	}
}

func TestFlushExists(t *testing.T) {
	// NOTE: not necessarily 5 cards, could be more
	const numTests = 4

	var theres = [numTests]CardSet{
		// with aces
		AceOfHearts | JackOfHearts | FiveOfHearts | FourOfHearts | TwoOfHearts,
		NineOfClubs | ThreeOfClubs | FourOfClubs | FiveOfClubs | JackOfClubs | SevenOfHearts | AceOfDiamonds,
		// without aces
		Diamonds & ^Aces,
		QueenOfDiamonds | KingOfDiamonds | SixOfDiamonds | TenOfDiamonds | TwoOfDiamonds | TwoOfClubs,
	}

	var expectedOutputs = [numTests]CardSet{
		AceOfHearts | JackOfHearts | FiveOfHearts | FourOfHearts | TwoOfHearts,
		NineOfClubs | ThreeOfClubs | FourOfClubs | FiveOfClubs | JackOfClubs,
		Diamonds & ^Aces,
		QueenOfDiamonds | KingOfDiamonds | SixOfDiamonds | TenOfDiamonds | TwoOfDiamonds,
	}

	for i := 0; i < numTests; i++ {
		var expectedOutput CardSet = expectedOutputs[i]
		var output CardSet = flush(theres[i])

		if output != expectedOutput {
			t.Errorf(errMsg(output, expectedOutput))
		}
	}
}

func TestFlushDoesNotExist(t *testing.T) {
	const numTests = 2

	var notTheres = [numTests]CardSet{
		// almost a flush (4)
		// good also because tests with aces (which is an important case)
		AceOfHearts | KingOfHearts | TenOfHearts | NineOfHearts | ThreeOfDiamonds | TwoOfClubs,
		// not almost a flush 3 or under
		QueenOfSpades | QueenOfHearts | QueenOfClubs,
	}

	var expectedOutput CardSet = 0

	for i := 0; i < numTests; i++ {
		var output CardSet = flush(notTheres[i])

		if output != expectedOutput {
			t.Errorf(errMsg(output, expectedOutput))
		}
	}
}

func TestFullHouseExists(t *testing.T) {
	const numTests = 4

	var theres = [numTests]CardSet{
		// exact
		AceOfSpades | AceOfHearts | AceOfDiamonds | TenOfClubs | TenOfDiamonds,
		FourOfClubs | FourOfDiamonds | FourOfHearts | FiveOfClubs | FiveOfDiamonds | SixOfSpades,
		// two trips
		(Sevens & ^SevenOfHearts) | (Kings & ^KingOfSpades),
		(Aces & ^AceOfSpades) | (Jacks & ^JackOfDiamonds) | SixOfHearts,
	}

	var expectedOutputs = [numTests]CardSet{
		AceOfSpades | AceOfHearts | AceOfDiamonds | TenOfClubs | TenOfDiamonds,
		FourOfClubs | FourOfDiamonds | FourOfHearts | FiveOfClubs | FiveOfDiamonds,
		(Sevens & ^SevenOfHearts) | (Kings & ^KingOfSpades),
		(Aces & ^AceOfSpades) | (Jacks & ^JackOfDiamonds),
	}

	for i := 0; i < numTests; i++ {
		var expectedOutput CardSet = expectedOutputs[i]
		var output CardSet = fullHouse(theres[i])

		if output != expectedOutput {
			t.Errorf(errMsg(output, expectedOutput))
		}
	}
}

func TestFullHouseDoesNotExist(t *testing.T) {

	const numTests = 6

	var notTheres = [numTests]CardSet{
		// test with triple
		AceOfSpades | AceOfClubs | AceOfHearts,
		TenOfHearts | TenOfSpades | TenOfClubs | TwoOfHearts | SevenOfDiamonds,
		// test with pair
		AceOfDiamonds | AceOfSpades,
		KingOfDiamonds | KingOfHearts | ThreeOfDiamonds | ThreeOfHearts | TwoOfDiamonds | TwoOfClubs,
		// test with nothing
		NoCards,
		// test with quad
		Fours,
	}

	var expectedOutput CardSet = 0

	for i := 0; i < numTests; i++ {
		var output CardSet = fullHouse(notTheres[i])

		if output != expectedOutput {
			t.Errorf(errMsg(output, expectedOutput))
		}
	}
}

// // pairs and trips are slightly different since we expect to find many, so we will use
// // function literals instead to test them

func TestTrips(t *testing.T) {
	const numTests = 9

	var inputs = [numTests]CardSet{
		// with 1 trip
		NineOfClubs | NineOfDiamonds | NineOfSpades,
		FourOfHearts | FourOfDiamonds | FourOfClubs | SevenOfDiamonds,
		// with 2 trips
		JackOfClubs | JackOfDiamonds | JackOfHearts | TenOfSpades | TenOfDiamonds | TenOfClubs | NineOfClubs,
		// with quad (should be ignored)
		Sevens,
		Threes | TwoOfHearts,
		// with no trips (with doubles)
		AceOfClubs | TenOfHearts | QueenOfSpades | KingOfClubs | FiveOfClubs | FiveOfDiamonds | TwoOfSpades,
		AceOfDiamonds | AceOfSpades,
		// with no trips (no doubles)
		Hearts,
		NoCards,
	}

	var isExpecteds = []func(trip1, trip2 CardSet) (bool, string) {
		// "with 1 trip"
		func(t1, t2 CardSet) (bool, string) {
			return (t1 == NineOfClubs | NineOfDiamonds | NineOfSpades) && t2 == 0, "9C 9D 9S, "
		},
		func(t1, t2 CardSet) (bool, string) {
			return (t1 == FourOfHearts | FourOfDiamonds | FourOfClubs) && t2 == 0, "4C 4D 4H, "
		},
		// "with 2 trips"
		func(t1, t2 CardSet) (bool, string) {
			return (
				t1 == JackOfClubs | JackOfDiamonds | JackOfHearts && 
				t2 == TenOfSpades | TenOfDiamonds | TenOfClubs), "JC JD JH, TC TD TS"
		},
		// "with quad (should be ignored)"
		func(t1, t2 CardSet) (bool, string) {
			return t1 | t2 == 0, ", , "
		},
		func(t1, t2 CardSet) (bool, string) {
			return t1 | t2 == 0, ", , "
		},
		// "with no trips (with doubles)"
		func(t1, t2 CardSet) (bool, string) {
			return t1 | t2 == 0, ", , "
		},
		func(t1, t2 CardSet) (bool, string) {
			return t1 | t2 == 0, ", , "
		},
		// "with no trips (no doubles)"
		func(t1, t2 CardSet) (bool, string) {
			return t1 | t2 == 0, ", , "
		},
		func(t1, t2 CardSet) (bool, string) {
			return t1 | t2 == 0, ", , "
		},
	}

	for i := 0; i < numTests; i++ {
		var output1, output2 CardSet = triplet(inputs[i])

		ok, message := isExpecteds[i](output1, output2)
		if !ok {
			var outputs string = fmt.Sprintf(
				"{ %s, %s }", 
				CardSetToString(output1),
				CardSetToString(output2),
			)

			t.Errorf(gotButExpected(outputs, message))
		}
	}
}

// in the future you will want to identify your tests
// maybe this is why google test was better...
func TestPairs(t *testing.T) {
	const numTests = 11

	var inputs = [numTests]CardSet{
		// try with 1 pair
		KingOfDiamonds | KingOfHearts,
		JackOfSpades | JackOfClubs | EightOfHearts | TenOfHearts | AceOfClubs | TwoOfHearts | ThreeOfSpades,
		// try with 2 pairs
		SixOfSpades | SixOfDiamonds | SevenOfHearts | SevenOfSpades | NineOfHearts | TenOfHearts,
		// try with 3 pairs
		AceOfClubs | AceOfDiamonds | KingOfHearts | KingOfClubs | ThreeOfSpades | ThreeOfHearts,
		SevenOfClubs | SevenOfDiamonds | SixOfHearts | SixOfClubs | EightOfSpades | EightOfHearts | QueenOfSpades,
		// try with triplets (should be ignored)
		AceOfClubs | AceOfHearts | AceOfSpades,
		TwoOfSpades | TwoOfClubs | TwoOfHearts | ThreeOfDiamonds | ThreeOfHearts,
		// try with quads (should be ignored)
		Kings,
		Aces | NineOfClubs | TenOfHearts,
		// try with no pairs
		NoCards,
		SixOfDiamonds | SevenOfSpades | EightOfSpades | NineOfClubs,
	}

	var isExpecteds = []func(pair1, pair2, pair3 CardSet) (bool, string) {
		// "try with 1 pair"
		func(p1, p2, p3 CardSet) (bool, string) {
			return (p1 == (KingOfDiamonds | KingOfHearts) && p2 | p3 == 0), "KD KH, , "
		},
		func(p1, p2, p3 CardSet) (bool, string) {
			return p1 == (JackOfSpades | JackOfClubs) && p2 | p3 == 0, "JS JC, , "
		},
		// "try with 2 pairs"
		func(p1, p2, p3 CardSet) (bool, string) {
			return (
				p1 == SevenOfSpades | SevenOfHearts && 
				p2 == SixOfSpades | SixOfDiamonds && 
				p3 == 0), "7S 7H, 6S 6D, "
		},
		// "try with 3 pairs"
		func(p1, p2, p3 CardSet) (bool, string) {
			return (
				p1 == AceOfDiamonds | AceOfClubs &&
				p2 == KingOfHearts | KingOfClubs &&
				p3 == ThreeOfSpades | ThreeOfHearts), "AD AC, KH KC, 3S 3H"
		},
		func(p1, p2, p3 CardSet) (bool, string) {
			return (
				p1 == EightOfSpades | EightOfHearts &&
				p2 == SevenOfClubs | SevenOfDiamonds &&
				p3 == SixOfHearts | SixOfClubs), "8S 8H, 7D 7C, 6H 6C"
		},
		// "try with triplets (should be ignored)"
		func(p1, p2, p3 CardSet) (bool, string) {
			return p1 | p2 | p3 == 0, ", , "
		},
		func(p1, p2, p3 CardSet) (bool, string) {
			return p1 == (ThreeOfDiamonds | ThreeOfHearts) && p2 | p3 == 0, "3H 3D, , "
		},
		// "try with quads (should be ignored)"
		func(p1, p2, p3 CardSet) (bool, string) {
			return p1 | p2 | p3 == 0, ", , "
		},
		func(p1, p2, p3 CardSet) (bool, string) {
			return p1 | p2 | p3 == 0, ", , "
		},
		// "try with no pairs"
		func(p1, p2, p3 CardSet) (bool, string) {
			return p1 | p2 | p3 == 0, "" //empty string means nothing
		},
		func(p1, p2, p3 CardSet) (bool, string) {
			return p1 | p2 | p3 == 0, "" //ibid
		},
	}

	for i := 0; i < numTests; i++ {
		var output1, output2, output3 CardSet = pair(inputs[i])
		
		ok, message := isExpecteds[i](output1, output2, output3)
		if !ok {
			var outputs string = fmt.Sprintf(
				"%s, %s, %s", 
				CardSetToString(output1),
				CardSetToString(output2),
				CardSetToString(output3),
			)

			t.Errorf(gotButExpected(outputs, message))
		}
	}
}

func TestHighCard(t *testing.T) {
	const numTests = 6

	var theres = [numTests]CardSet{
		// try singletons
		AceOfDiamonds,
		ThreeOfClubs,
		// try a pair
		SixOfHearts | SixOfSpades,
		// try sets
		AllCards,
		SevenOfDiamonds | SixOfSpades | FourOfSpades,
		NoCards,
	}

	var expectedOutputs = [numTests]CardSet{
		AceOfDiamonds,
		ThreeOfClubs,
		SixOfSpades,
		AceOfSpades,
		SevenOfDiamonds,
		NoCards,
	}

	for i := 0; i < numTests; i++ {
		var expectedOutput CardSet = expectedOutputs[i]
		var output CardSet = extractHighCard(theres[i])

		if output != expectedOutput {
			t.Errorf(errMsg(output, expectedOutput))
		}
	}
}


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