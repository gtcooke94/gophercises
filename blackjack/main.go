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
	printPlayers((*gamePtr).Players)
	takeTurns(gamePtr)
	finalScoring(gamePtr)
}

func finalScoring(g *blackjack.Game) {
	printPlayers((*g).Players)
}

func takeTurns(g *blackjack.Game) {
	for i := range (*g).Players {
		playerPtr := &((*g).Players[i])
		continueTurn := true
		if playerPtr.DealerFlag {
			fmt.Println("====================")
			(*g).DealerTurn(playerPtr)
			continueTurn = false
		}
		for continueTurn {
			fmt.Println("====================")
			fmt.Println(*playerPtr)
			printOptions()
			action, err := acceptOption()
			for err != nil {
				fmt.Println("Invalid option selected. Please selection a valid option")
				printOptions()
				action, err = acceptOption()

			}
			continueTurn = (*g).PlayerTurn(playerPtr, blackjack.Action(action))
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
