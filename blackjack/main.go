package main

import (
	"bufio"
	"fmt"
	"github.com/gtcooke94/gophercises/blackjack/game"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("vim-go")
	gamePtr := blackjack.StartGame(2, 2)
	play := true
	for play {
		startRound(gamePtr)
		printPlayers((*gamePtr).Players)
		takeTurns(gamePtr)
		finalScoring(gamePtr)
		cleanupRound(gamePtr)
		play = playAgain()
	}
}

func startRound(g *blackjack.Game) {
	fmt.Println("==================== Starting Round... ====================")
	g.StartRound()
}

func cleanupRound(g *blackjack.Game) {
	fmt.Println("==================== Cleaning up... ====================")
}

func finalScoring(g *blackjack.Game) {
	printPlayers((*g).Players)
}

func playAgain() bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Play again?\n0. No\n1. Yes")
	play, _ := reader.ReadString('\n')
	play = strings.Replace(play, "\n", "", -1)
	if play == "1" {
		return true
	}
	return false
}

func takeTurns(g *blackjack.Game) {
	for i := range (*g).Players {
		playerPtr := &((*g).Players[i])
		continueTurn := true
		if playerPtr.DealerFlag {
			fmt.Println("====================")
			(*g).DealerTurn(playerPtr)
			// continueTurn = false
			continue
		}
		fmt.Printf("==================== Player %d's Turn ====================\n", i)
		fmt.Println(*playerPtr)
		for continueTurn {
			printOptions()
			action, err := acceptOption()
			for err != nil {
				fmt.Println("Invalid option selected. Please selection a valid option")
				printOptions()
				action, err = acceptOption()

			}
			continueTurn = (*g).PlayerTurn(playerPtr, blackjack.Action(action))
			fmt.Println(*playerPtr)
		}
	}
}

func acceptOption() (int, error) {
	reader := bufio.NewReader(os.Stdin)
	actionStr, _ := reader.ReadString('\n')
	actionStr = strings.Replace(actionStr, "\n", "", -1)
	action, err := strconv.Atoi(actionStr)
	return action, err
}

func printOptions() {
	for i, action := range blackjack.Actions {
		fmt.Printf("%d. %s\n", i, action)
	}
}

func printPlayers(players []blackjack.Player) {
	for _, player := range players {
		fmt.Println(player)
	}
}
