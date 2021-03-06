package poker

import (
	"fmt"
	"math/bits"
	"strings"
)

// 13 * 4 = 52 unique cards in poker
// any set of cards can be fit in a >= 52 bit vector
// we'll use uint64 for sets of cards
// we'll encode Ace as both 1s and 14s for now for convenience...

// this will leave us with 8 new possible cards for future expansions
// which corrresponds to 2 new possible cards per suit if we so desire
// (i.e. you may want to add jokers)

// we will define each block of four bits to correspond to the number
// and they'll be ordered by their numerical value
// and they'll be ordered by ascending suit

type CardSet uint64

// these are singleton sets (i.e. of single cards)
const (
	// small aces (ones)
	_Ace1ofClubs CardSet = 1 << iota
	_Ace1ofDiamonds
	_Ace1ofHearts
	_Ace1ofSpades
	// twos
	TwoOfClubs
	TwoOfDiamonds
	TwoOfHearts
	TwoOfSpades
	// threes
	ThreeOfClubs
	ThreeOfDiamonds
	ThreeOfHearts
	ThreeOfSpades
	// fours
	FourOfClubs
	FourOfDiamonds
	FourOfHearts
	FourOfSpades
	// fives
	FiveOfClubs
	FiveOfDiamonds
	FiveOfHearts
	FiveOfSpades
	// sixes
	SixOfClubs
	SixOfDiamonds
	SixOfHearts
	SixOfSpades
	// sevens
	SevenOfClubs
	SevenOfDiamonds
	SevenOfHearts
	SevenOfSpades
	// eights
	EightOfClubs
	EightOfDiamonds
	EightOfHearts
	EightOfSpades
	// nines
	NineOfClubs
	NineOfDiamonds
	NineOfHearts
	NineOfSpades
	// tens
	TenOfClubs
	TenOfDiamonds
	TenOfHearts
	TenOfSpades
	// jacks (elevens)
	JackOfClubs
	JackOfDiamonds
	JackOfHearts
	JackOfSpades
	// queens (twelves)
	QueenOfClubs
	QueenOfDiamonds
	QueenOfHearts
	QueenOfSpades
	// kings (thirteens)
	KingOfClubs
	KingOfDiamonds
	KingOfHearts
	KingOfSpades
	// big aces (fourteens)
	_Ace2ofClubs
	_Ace2ofDiamonds
	_Ace2ofHearts
	_Ace2ofSpades
)

// for convenience's sake we will treat Ace as two cards
// which is analogous to having the option to use it as
// a "little" card or as a high card
// (useful for computing straights)
const (
	AceOfClubs CardSet = (_Ace1ofClubs | _Ace2ofClubs) << iota
	AceOfDiamonds
	AceOfHearts
	AceOfSpades
)

// used to check for straights, pairs, triplets, full houses, and four of a kind
const (
	_Ace1s CardSet = (_Ace1ofClubs | _Ace1ofDiamonds | _Ace1ofHearts | _Ace1ofSpades) << (4 * iota)
	Twos
	Threes
	Fours
	Fives
	Sixes
	Sevens
	Eights
	Nines
	Tens
	Jacks
	Queens
	Kings
	_Ace2s
)

const Aces CardSet = _Ace1s | _Ace2s

// used to check for flushes
const (
	Clubs CardSet = (AceOfClubs | TwoOfClubs | ThreeOfClubs | FourOfClubs |
		FiveOfClubs | SixOfClubs | SevenOfClubs | EightOfClubs |
		NineOfClubs | TenOfClubs | JackOfClubs | QueenOfClubs |
		KingOfClubs) << iota
	Diamonds
	Hearts
	Spades
)

const (
	ClubsRoyalFlush CardSet = (AceOfClubs | KingOfClubs | QueenOfClubs | JackOfClubs | TenOfClubs) << iota
	DiamondsRoyalFlush
	HeartsRoyalFlush
	SpadesRoyalFlush
)

const AllCards = Clubs | Diamonds | Hearts | Spades
const NoCards = 0

// to visualize use
// fmt.Printf("% 064b", n) for a CardSet n
// I'll probably use this in tests
// note the 0 means to keep leading zeros
// before tests play around with https://play.golang.org/

// return the royal flush if it's a royal flush else zero
func royalFlush(cardset CardSet) CardSet {
	var royalFlush CardSet = SpadesRoyalFlush
	var mask CardSet = royalFlush & cardset
	for royalFlush > ClubsRoyalFlush-1 {
		if mask == royalFlush {
			return mask
		}

		royalFlush >>= 1
		mask = royalFlush & cardset
	}
	return 0
}

// return 0 if there is not four of a kind
// otherwise return the highest four of a kind they hold
// as a card set (note, it's only possible to have one four of a kind)
func fourOfAKind(cardset CardSet) CardSet {
	if cardset&Aces == Aces {
		return Aces
	}

	var quad CardSet = Kings
	for quad > _Ace1s {
		if quad&cardset == quad {
			return quad
		}

		quad >>= 4
	}

	return 0
}

// return the highest set of cards that form the straight flush if there is one
// (note: the highest set of flush and straight cards may not be the same)
func straightFlush(cardset CardSet) CardSet {
	var flush CardSet = Spades
	var straightSet = straight(cardset & flush)
	for straightSet == 0 && flush >= Clubs {
		flush >>= 1
		straightSet = straight(cardset & flush)
	}

	return straightSet
}

// return the highest full house
// NOTE that pair and triplet are already exclusive of eachother
// NOTE that it CANNOT catch a fullhouse when there is quads (due to the construction of pair and trips)
// that's OK because the quads beat the fullhouse
func fullHouse(cardset CardSet) CardSet {
	p1, _, _ := pair(cardset)
	t1, t2 := triplet(cardset)

	if p1 > 0 && t1 > 0 {
		return p1 | t1
	} else if t1 > 0 && t2 > 0 {
		return t1 | t2 // not perfectly efficient but ok for now
	} else {
		return 0
	}
}

// return the highest set of cards that form the flush if there is one or 0
func flush(cardset CardSet) CardSet {
	var flush CardSet = Spades
	var mask CardSet = flush & cardset
	for flush > Clubs-1 {
		// this is really hacky and pretty cool, check it out
		bc := bits.OnesCount64(uint64(mask))
		// note that if you have an ace you get 1 more bit count so we should check if there is an ace
		// and then do 1 more if there is potentially an ace
		if (Aces&mask == 0 && bc >= 5) || bc >= 6 {
			return mask
		}

		flush >>= 1
		mask = flush & cardset
	}
	return 0
}

// return the highest set of cards that form the straight (if there is one) or 0
// (note: pairs of cards of the same value WILL be returned - feature not bug!)
func straight(cardset CardSet) CardSet {
	// strategy is to keep track of one card above the highest card in the straight (highCard)
	// and as soon as we break the straight with an insufficint count, clear the upper
	// section of the cardset (passed by value) so we can keep only the straight inside the output
	// use lowCard to keep track of the next possible highCard and clear the lower bits as well

	var quad CardSet = Aces
	// NOTE: this will NOT work after we add additional cards
	var highCard CardSet = _Ace2ofSpades << 1
	var lowCard CardSet = _Ace2ofClubs // same as KingOfSpades << 1
	var count uint32 = 0

	for quad > 0 {
		if quad&cardset > 0 {
			if count == 4 {
				if cardset&_Ace2s > 0 {
					// deal with upper ace (need to add back in lower ace)
					// this ONLY happens if we found a flush with ace high
					return (cardset & ^(lowCard - 1)) | ((cardset & _Ace2s) >> 52)
				} else if cardset&_Ace1s > 0 && lowCard == _Ace1ofClubs {
					// need to deal with lower Ace (deleted upper ace need to add back in)
					// this ONLY happens if we have ace bottom (i.e. 5 4 3 2 Ace)
					return cardset | ((cardset & _Ace1s) << 52)
				} else {
					// note lowCard is at LEAST 1 (i.e. _Ace1ofClubs)
					// but when it hits zero then the quad will hit zero
					// note tha the quad will be on _Ace1s when this is the lower card
					return cardset & ^(lowCard - 1)
				}
			}

			count += 1
		} else {
			count = 0
			cardset &= (highCard - 1)
			highCard = lowCard
		}

		lowCard >>= 4
		quad >>= 4
	}

	return 0
}

func maskIsTriplet(mask CardSet, quad CardSet, highestCard CardSet) bool {
	return (mask == quad & ^highestCard ||
		mask == quad & ^(highestCard>>1) ||
		mask == quad & ^(highestCard>>2) ||
		mask == quad & ^(highestCard>>3))
}

// return 0, 0 if they have no triplets
// otherwise return highest_triplet, 0 | lowest_triplet
// (note that there are at most two triplets)
func triplet(cardset CardSet) (CardSet, CardSet) {
	var quad CardSet = Aces
	var highestCard CardSet = AceOfSpades
	var highTriplet CardSet = 0
	var lowTriplet CardSet = 0

	for quad > _Ace1s {
		var mask CardSet = quad & cardset

		// it's one of three posible triplets (gotten by having a quad
		// and losing one of four possible single cards)
		var isTriplet bool = maskIsTriplet(mask, quad, highestCard)

		// can be made branchless with maxes/mins potentially
		if isTriplet && highTriplet > 0 {
			lowTriplet = mask
			break
		} else if isTriplet {
			highTriplet = mask
		}

		quad >>= 4
		highestCard >>= 4
	}

	return highTriplet, lowTriplet
}

// return 0, 0, 0 if they have no pairs
// otherwise return highest pair, 0 | second highest pair, 0 | lowest pair
// (note can have at most three pairs)
// (note: ignores quads and ignores triplets, those won't overlap)
func pair(cardset CardSet) (CardSet, CardSet, CardSet) {
	var quad CardSet = Aces
	var highestCard CardSet = AceOfSpades

	var highPair CardSet = 0
	var medPair CardSet = 0
	var lowPair CardSet = 0

	for quad > _Ace1s {
		var mask CardSet = quad & cardset

		// this is a pair if it's not a quad, triplet, or singleton
		// could also use bitcount (and don't forget to check it's nonzero, LOL)
		if (mask > 0) && (mask != quad) && !maskIsTriplet(mask, quad, highestCard) && singleton(mask) == 0 {
			// consider whether to make branchless? idea is maxes and mins
			if highPair == 0 {
				highPair = mask
			} else if medPair == 0 {
				medPair = mask
			} else {
				lowPair = mask
				break
			}
		}

		quad >>= 4
		highestCard >>= 4
	}

	return highPair, medPair, lowPair
}

// return card if it's a singleton else 0
func singleton(cardset CardSet) CardSet {
	// helpful: http://blog.marcinchwedczuk.pl/how-to-check-if-a-number-is-a-power-of-two

	// a number is a single card IFF it's a power of 2 or it's a single ace
	if cardset == AceOfSpades || cardset == AceOfHearts || cardset == AceOfDiamonds || cardset == AceOfClubs {
		return cardset
	} else if (cardset > _Ace1s) && (cardset > 0) && (cardset&(cardset-1)) == 0 {
		return cardset
	}

	return 0
}

// return the highest card or zero if there are no cards
func extractHighCard(cardset CardSet) CardSet {
	var test CardSet = AceOfSpades

	// NOTE: when you bitshift something that goes too far right its clipped
	// so it's ok to start with aces!
	for (test&cardset != test) && test > _Ace1ofSpades {
		test >>= 1
	}

	if test > _Ace1ofSpades {
		return test
	} else {
		return 0
	}
}

// no need to test this, look how simple it is
// I've tested it manually
// use it to debug
func bitvectorString(cardset CardSet) string {
	return fmt.Sprintf("% 064b", cardset)
}

// return a string of <number><suit> for example Ace of Hearts -> AH, 10 of Clubs -> 10C
// and it will be space separated and sorted from min to max
func CardSetToString(cardset CardSet) string {
	// we look at the set card by card and we start
	// by looking at the lowest card, the two of clubs
	var card CardSet = TwoOfClubs

	// this is the number of shifts from 0
	// i.e. 1 = 1 << 0, 2 = 1 << 1, 4 = 1 << 2, 8 = 1 << 3, so 0, 1, 2, 3, ...
	var shift byte = 4

	// Recall there are four extra aces...
	// this will basically add an offset to the shift index / 4 (to get a step function
	// corresponding to the number, due to integer division) and then use different offsets for
	// different types (actually only numbers need offsets, and rest can use lookup tables, but
	// this is good enough for our purposes and simple to write)

	// ascii: (cards) 0:48, A: 65, K: 75, Q: 81, J: 74
	var shiftBounds = [6]byte{
		36, // numbers
		40, // ten
		44, // jack
		48, // queen
		52, // king
		56, // ace
	}
	var asciiOffsets = [6]byte{
		49, // numbers
		75, // ten: down 9 because of shift / 4 is up by 9
		64, // jack: down by 10 because shift / 4 is up by 10
		70, // queen: down by 11 (ibid)
		63, // king: down by 12 (ibid)
		52, // ace: down by 13 (ibid)
		// those shifts were zero indexxed
	}

	// ascii: (suits) S: 83, H: 72, D: 68, C: 67
	// lookup by modulus (i.e. suit is clubs mod 0, diamonds mod 1 etc...)
	var suitLookup = [4]byte{67, 68, 72, 83}

	var b strings.Builder
	for i := 0; i < 6; i++ {
		for shift < shiftBounds[i] {
			if card&cardset > 0 {
				b.WriteByte(shift/4 + asciiOffsets[i])
				b.WriteByte(suitLookup[shift%4])
				b.WriteByte(' ')
			}

			shift++
			card <<= 1
		}
	}

	return strings.TrimRight(b.String(), " ")
}
