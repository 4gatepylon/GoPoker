package poker

type Hand interface {
	Beats(other Hand) bool
}

type CardSetHand struct {
	cards CardSet
}

func (this CardSetHand) Beats(other Hand) {
	switch other.(type) {
	case CardSetHand:
		return firstBeats(this, other.(CardSetHand))
	default:
		panic("Not implemented!")
	}
}

func firstBeats(beater CardSetHand, loser CardSetHand) {
	var beaterCards CardSet = beater.cards
	var loserCards CardSet = loser.cards
	return false
}