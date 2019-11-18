//go:generate stringer -type=Action
package blackjack

import (
	"fmt"
	"github.com/gtcooke94/gophercises/deck"
	"math"
)

type Player struct {
	Cards            []deck.Card
	Name             string
	Chips            int
	DealerFlag       bool
	DealerHiddenFlag bool
}

type Action uint8

const MAXSCORE = 21

const (
	Hit Action = iota
	Pass
)

var Actions = [...]Action{Hit, Pass}

type Game struct {
	Players  []Player
	gameDeck []deck.Card
}

// var gameDeck []deck.Card = make([]deck.Card, 0)
// var players []Player = make([]Player, 0)
// var gameDeckPtr *[]deck.Card = &gameDeck
// var playersPtr *[]Player = &players

var game Game

func (g *Game) drawCard() deck.Card {
	toReturn := (*g).gameDeck[len((*g).gameDeck)-1]
	(*g).gameDeck = (*g).gameDeck[:len((*g).gameDeck)-1]
	return toReturn
}

func (g *Game) addPlayer() {
	card1 := g.drawCard()
	card2 := g.drawCard()
	newPlayer := newPlayer([]deck.Card{card1, card2}, fmt.Sprintf("Player %d", len(g.Players)))
	(*g).Players = append((*g).Players, newPlayer)
}

func (g *Game) addDealer() {
	card1 := g.drawCard()
	card2 := g.drawCard()
	(*g).Players = append((*g).Players, newDealer([]deck.Card{card1, card2}))
}

func StartGame(nPlayers int, nDecks int) *Game {
	fmt.Println("Starting Game...")
	gameDeck := deck.New(deck.Deck(nDecks), deck.Shuffle)
	game = Game{make([]Player, 0), gameDeck}
	for i := 0; i < nPlayers; i++ {
		game.addPlayer()
	}
	game.addDealer()
	return &game
}

func (p Player) Score() int {
	score, _ := scoreLogic(p.Cards)
	return score
}

// Returns score, flag indicating if an ace counted as 11 or not
func scoreLogic(cards []deck.Card) (int, bool) {
	numAces := 0
	score := 0
	aceAs11Flag := false
	for _, card := range cards {
		if card.Rank == deck.Ace {
			numAces++
		} else if card.Rank >= deck.Ten {
			score = score + 10
		} else {
			score = score + int(card.Rank)
		}
	}
	// Score >= 11 all aces count as 1
	if score >= 11 {
		score = score + numAces
	} else if numAces >= 1 {
		if score+11+numAces-1 > MAXSCORE {
			score = score + numAces
		} else {
			score = score + 11 + numAces - 1
			aceAs11Flag = true
		}
	}
	return score, aceAs11Flag
}

func (p Player) String() string {
	var ret string
	fmt.Println("DEBUG", p.DealerHiddenFlag, p.DealerFlag)
	if p.DealerHiddenFlag && p.DealerFlag {
		ret = fmt.Sprintf("%s: Cards: %s, HIDDEN\n", p.Name, p.Cards[0])
	} else {
		ret = fmt.Sprintf("%s: Score: %d. Cards: %v\n", p.Name, p.Score(), p.Cards)
	}
	return ret
}

func (g *Game) PlayerTurn(p *Player, action Action) bool {
	if (*p).Score() > MAXSCORE {
		return false
	}
	switch action {
	case Hit:
		newCard := g.drawCard()
		(*p).Cards = append((*p).Cards, newCard)
		return true
	case Pass:
		return false
	}
	return false
}

func (g *Game) DealerTurn(p *Player) {
	// Dealer hits on <=16, soft 17, holds otherwise
	for hit := dealerHitCondition(*p); hit; {
		newCard := g.drawCard()
		(*p).Cards = append((*p).Cards, newCard)
		hit = dealerHitCondition(*p)
	}
	(*p).DealerHiddenFlag = false
	fmt.Println("DEBUG", (*p).DealerHiddenFlag)
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

func newPlayer(c []deck.Card, name string) Player {
	return Player{Cards: c, Name: name, Chips: 0, DealerFlag: false, DealerHiddenFlag: false}
}

func newDealer(c []deck.Card) Player {
	return Player{Cards: c, Name: "Dealer", Chips: int(math.Inf(1)), DealerFlag: true, DealerHiddenFlag: true}
}
