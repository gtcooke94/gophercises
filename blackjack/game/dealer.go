package blackjack

import (
	"github.com/gtcooke94/gophercises/deck"
	"math"
)

func (g *Game) DealerTurn(p *Player) {
	// Dealer hits on <=16, soft 17, holds otherwise
	for hit := dealerHitCondition(*p); hit; {
		newCard := g.drawCard()
		(*p).Cards = append((*p).Cards, newCard)
		hit = dealerHitCondition(*p)
	}
	(*p).DealerHiddenFlag = false
}

func dealerHitCondition(p Player) bool {
	if p.Score() <= 16 {
		return true
	} else if p.Score() == 17 {
		_, soft17 := scoreLogic(p.Cards)
		if soft17 {
			return true
		}
	}
	return false

}

func newDealer() Player {
	return Player{Cards: make([]deck.Card, 0), Name: "Dealer", Chips: int(math.Inf(1)), DealerFlag: true, DealerHiddenFlag: true, CurrentBet: 0}
}
