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
	CurrentBet       int
}

type Action uint8

const MAXSCORE = 21

const MINBET = 1

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

func (g *Game) DetermineWinners() ([]*Player, []*Player, []*Player) {
	winners := make([]*Player, 0)
	ties := make([]*Player, 0)
	losers := make([]*Player, 0)
	dealerScore := g.Players[len(g.Players)-1].Score()
	for i := range g.Players {
		player := &g.Players[i]
		if player.DealerFlag {
			continue
		}
		playerScore := player.Score()
		// If a player busts, they lose
		if playerScore > MAXSCORE {
			losers = append(losers, player)
		} else if dealerScore > MAXSCORE {
			// If the dealer busts and the player didn't, the player wins
			winners = append(winners, player)
		} else {
			// Neither busted, compare scores
			if playerScore > dealerScore {
				winners = append(winners, player)
			} else if playerScore == dealerScore {
				ties = append(ties, player)
			} else {
				losers = append(losers, player)
			}
		}
	}
	return winners, ties, losers
}

func (g *Game) PayoutWin(p *Player) {
	// TODO after implementing betting
	p.Chips = p.Chips + 2*p.CurrentBet
}

func (g *Game) TakeBet(p *Player) {
	// TODO after implementing betting
}

func (g *Game) TiedDealer(p *Player) {
	// TODO after implementing betting
	p.Chips = p.Chips + p.CurrentBet
}

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

func (g *Game) PlaceBet(player *Player, bet int) bool {
	if bet < MINBET {
		return false
	}
	if bet > player.Chips {
		return false
	}
	player.CurrentBet = bet
	player.Chips = player.Chips - bet
	return true
}

func (g *Game) DealCards() {
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
		player.CurrentBet = 0
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
	} else if p.DealerFlag {
		ret = fmt.Sprintf("%s: Score: %s.\nCards: %v\n", p.Name, scoreStr, p.Cards)
	} else {
		ret = fmt.Sprintf("%s: Score: %s.\nCards: %v\nChips: %d\nCurrent Bet: %d", p.Name, scoreStr, p.Cards, p.Chips, p.CurrentBet)
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
	return Player{Cards: make([]deck.Card, 0), Name: name, Chips: 100, DealerFlag: false, DealerHiddenFlag: false, CurrentBet: 0}
}
