//go:generate stringer -type=Action
package blackjack

import (
	"fmt"
	"github.com/gtcooke94/gophercises/deck"
	"strconv"
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
	Stand
)

var Actions = [...]Action{Hit, Stand}

type Game struct {
	Players  []Player
	gameDeck []deck.Card
	discard  []deck.Card
}

// var gameDeck []deck.Card = make([]deck.Card, 0)
// var players []Player = make([]Player, 0)
// var gameDeckPtr *[]deck.Card = &gameDeck
// var playersPtr *[]Player = &players

var game Game

func (g *Game) drawCard() deck.Card {
	if len(g.gameDeck) == 0 {
		fmt.Println("The deck is out of cards... Shuffling up discarded card")
		g.gameDeck = deck.Shuffle(g.discard)
	}
	toReturn := (*g).gameDeck[len((*g).gameDeck)-1]
	(*g).gameDeck = (*g).gameDeck[:len((*g).gameDeck)-1]
	return toReturn
}

func (g *Game) addPlayer() {
	newPlayer := newPlayer(fmt.Sprintf("Player %d", len(g.Players)))
	(*g).Players = append((*g).Players, newPlayer)
}

func (g *Game) addDealer() {
	(*g).Players = append((*g).Players, newDealer())
}

func (g *Game) StartRound() {
	for i := range g.Players {
		player := &(g.Players[i])
		card1 := g.drawCard()
		card2 := g.drawCard()
		player.Cards = []deck.Card{card1, card2}
	}
}
func StartGame(nPlayers int, nDecks int) *Game {
	fmt.Println("Starting Game...")
	gameDeck := deck.New(deck.Deck(nDecks), deck.Shuffle)
	game = Game{make([]Player, 0), gameDeck, make([]deck.Card, 0)}
	for i := 0; i < nPlayers; i++ {
		game.addPlayer()
	}
	game.addDealer()
	return &game
}

func (g *Game) CleanupRound() {
	for i := range g.Players {
		player := g.Players[i]
		for _, c := range player.Cards {
			g.discard = append(g.discard, c)
		}
		player.Cards = make([]deck.Card, 0)
	}
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
	score := p.Score()
	scoreStr := strconv.Itoa(score)
	if score > MAXSCORE {
		scoreStr = "BUST"
	}
	if p.DealerHiddenFlag && p.DealerFlag {
		ret = fmt.Sprintf("%s: Cards: %s, HIDDEN\n", p.Name, p.Cards[0])
	} else {
		ret = fmt.Sprintf("%s: Score: %s. Cards: %v\n", p.Name, scoreStr, p.Cards)
	}
	return ret
}

func (g *Game) PlayerTurn(p *Player, action Action) bool {
	switch action {
	case Hit:
		newCard := g.drawCard()
		(*p).Cards = append((*p).Cards, newCard)
		if (*p).Score() >= MAXSCORE {
			return false
		}
		return true
	case Stand:
		return false
	}
	return false
}

func newPlayer(name string) Player {
	return Player{Cards: make([]deck.Card, 0), Name: name, Chips: 0, DealerFlag: false, DealerHiddenFlag: false}
}
