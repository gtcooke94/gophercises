//go:generate stringer -type=Suit,Rank

package deck

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

func main() {
	fmt.Println("vim-go")
}

// Card represents a playing card in a deck of cards.
type Card struct {
	Suit
	Rank
}

// Suit : The suit of a playing card
type Suit uint8

// iota starts at 0 and increments by one, keeping the same type
const (
	Spade Suit = iota
	Diamond
	Club
	Heart
	Joker
)

var suits = [...]Suit{Spade, Diamond, Club, Heart}

// Rank represents the value of the card (2-10, Jack, Queen, King, Ace)
type Rank uint8

// Having the first be an underscore makes values line up nicely
const (
	_ Rank = iota
	Ace
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

const (
	minRank = Ace
	maxRank = King
)

func (c Card) String() string {
	if c.Suit == Joker {
		return c.Suit.String()
	}
	return fmt.Sprintf("%s of %ss", c.Rank.String(), c.Suit.String())
}

// New - builds a new deck of cards
// Uses Functional Options!
// Takes variables amount of functions that take in a []Card and return a []Card
func New(opts ...func([]Card) []Card) []Card {
	var cards []Card
	// How can we loop through our constants?
	// Make a suits list because its short, then for the other loop over ranks
	for _, suit := range suits {
		for rank := minRank; rank <= maxRank; rank++ {
			cards = append(cards, Card{suit, rank})
		}
	}
	for _, opt := range opts {
		cards = opt(cards)
	}
	return cards
}

// Sorting functional options:
// Need a way to compare them
func absRank(c Card) int {
	return int(c.Suit)*int(maxRank) + int(c.Rank)
}

// Less : For sort, it takes in a slice and a less function
func Less(cards []Card) func(i, j int) bool {
	return func(i, j int) bool {
		return absRank(cards[i]) < absRank(cards[j])
	}
}

// DefaultSort functional option to sort by default
func DefaultSort(cards []Card) []Card {
	sort.Slice(cards, Less(cards))
	return cards
}

// Sort : Defining a custom sort function
// Input: a less funcion that takes in a slice of cards and returns a function that takes in two indices and returns if card at index i is less than card at index j.
// Output: function that takes in a slice of cards and returns a slice of cards
func Sort(less func(card []Card) func(i, j int) bool) func([]Card) []Card {
	return func(cards []Card) []Card {
		sort.Slice(cards, less(cards))
		return cards
	}
}

var shuffleRand = rand.New(rand.NewSource(time.Now().Unix()))

// // Shuffle : Shuffles the cards
// Original Shuffle
// func Shuffle(cards []Card) []Card {
//     rand.Seed(time.Now().UnixNano())
//     rand.Shuffle(len(cards), func(i, j int) { cards[i], cards[j] = cards[j], cards[i] })
//     return cards
// }

// Shuffle : the implementation from Gophercises
func Shuffle(cards []Card) []Card {
	ret := make([]Card, len(cards))
	// Basically gives you everything in the `random` package but with a different seed
	perm := shuffleRand.Perm(len(cards))
	for i, j := range perm {
		ret[i] = cards[j]
	}
	return ret
}

// Jokers : Place n jokers in the deck
func Jokers(n int) func([]Card) []Card {
	return func(cards []Card) []Card {
		for i := 0; i < n; i++ {
			cards = append(cards, Card{Suit: Joker, Rank: Rank(i)})
		}
		return cards
	}
}

// Filter : Filters cards based on the filter function provided
// The provided filter returns True if you want to remove the card, false otherwise
func Filter(f func(card Card) bool) func([]Card) []Card {
	return func(cards []Card) []Card {
		var ret []Card
		for _, card := range cards {
			if !f(card) {
				ret = append(ret, card)
			}
		}
		return ret
	}
}

// Deck : Create n copies of the deck
func Deck(n int) func([]Card) []Card {
	return func(cards []Card) []Card {
		var ret []Card
		for i := 0; i < n; i++ {
			// Take everything in cards and add it to ret
			ret = append(ret, cards...)
		}
		return ret
	}
}
