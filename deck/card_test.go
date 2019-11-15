package deck

import (
	"fmt"
	"math/rand"
	"testing"
)

func ExampleCard() {
	fmt.Println(Card{Rank: Ace, Suit: Heart})
	fmt.Println(Card{Rank: Two, Suit: Spade})
	fmt.Println(Card{Rank: Eight, Suit: Diamond})
	fmt.Println(Card{Rank: King, Suit: Club})
	fmt.Println(Card{Suit: Joker})

	// This will check the outputs when you run go test!
	// Output:
	// Ace of Hearts
	// Two of Spades
	// Eight of Diamonds
	// King of Clubs
	// Joker
}

func TestNew(t *testing.T) {
	cards := New()
	if len(cards) != 13*4 {
		t.Error("Wrong number of cards")
	}
}

func TestDefaultSort(t *testing.T) {
	cards := New(DefaultSort)
	if cards[0].Suit != Spade || cards[0].Rank != Ace {
		t.Error("Deck not sorted properly. Received:", cards[0])
	}
}

func TestSort(t *testing.T) {
	cards := New(Sort(Less))
	if cards[0].Suit != Spade || cards[0].Rank != Ace {
		t.Error("Deck not sorted properly. Received:", cards[0])
	}
}

// Bad test because it can fail
func TestShuffleOption(t *testing.T) {
	cards := New(DefaultSort, Shuffle)
	// cards := New(Shuffle, DefaultSort)
	if cards[0].Suit == Spade && cards[0].Rank == Ace {
		t.Error(fmt.Sprintf("Shuffled but %s is still at the bottom. Check the shuffle method", cards[0].String()))
	}
}

func TestShuffle(t *testing.T) {
	shuffleRand = rand.New(rand.NewSource(0))
	// We know that [40, 35] will be the first two with this seed
	unshuffled := New()
	shuffled := New(Shuffle)
	if unshuffled[40] != shuffled[0] && unshuffled[35] != shuffled[1] {
		t.Error("Something is wrong with shuffling")
	}
}

func TestJokers(t *testing.T) {
	n := 3
	cards := New(Jokers(n))
	numJokers := 0
	for _, c := range cards {
		if c.Suit == Joker {
			numJokers++
		}
	}
	if numJokers != n {
		t.Error("Not the right number of jokers")
	}
}

func TestFilter(t *testing.T) {
	filterFunc := func(c Card) bool {
		return c.Rank == 2 || c.Rank == 3
	}
	cards := New(Filter(filterFunc))
	for _, card := range cards {
		if card.Rank == Two || card.Rank == Three {
			t.Error("Deck filter did not work")
		}
	}

}

func TestDeck(t *testing.T) {
	numDecks := 4
	cards := New(Deck(numDecks))
	if len(cards) != 13*4*numDecks {
		t.Error("Wrong number of total cards")
	}
}
