package blackjack

import (
	"fmt"
	"github.com/gtcooke94/gophercises/deck"
	"math"
)

type Player struct {
	Cards      []deck.Card
	Name       string
	Chips      int
	DealerFlag bool
}

type Action uint8

const (
	Hit Action = iota
	Pass
)

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
	numAces := 0
	score := 0
	for _, card := range p.Cards {
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
		if score+11+numAces-1 > 21 {
			score = score + numAces
		} else {
			score = score + 11 + numAces - 1
		}
	}
	return score
}

func newPlayer(c []deck.Card, name string) Player {
	return Player{Cards: c, Name: name, Chips: 0, DealerFlag: false}
}

func newDealer(c []deck.Card) Player {
	return Player{Cards: c, Name: "Dealer", Chips: int(math.Inf(1)), DealerFlag: true}
}
