package main

import (
	"fmt"
	"github.com/gtcooke94/gophercises/blackjack/game"
)

func main() {
	fmt.Println("vim-go")
	players := blackjack.StartGame(2, 2)
	printPlayers(players)
}

func printPlayer(p blackjack.Player) {
	if !p.DealerFlag {
		fmt.Printf("%s: Score: %d. Cards: %v\n", p.Name, p.Score(), p.Cards)
	} else {
		fmt.Printf("%s: Cards: %s, HIDDEN\n", p.Name, p.Cards[0])
	}
}

func printPlayers(players *[]blackjack.Player) {
	for _, player := range *players {
		printPlayer(player)
	}
}
