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
	fmt.Println("==================== Place Bets... ====================")
	for i := range g.Players {
		player := &(g.Players[i])
		if player.DealerFlag {
			continue
		}
		fmt.Printf("%s Place Your Bet!\nChips: %d\n", player.Name, player.Chips)
		acceptBet(g, player)
	}
	fmt.Println("==================== Dealing Cards... ====================")
	g.DealCards()
}

func cleanupRound(g *blackjack.Game) {
	fmt.Println("==================== Cleaning up... ====================")
}

func finalScoring(g *blackjack.Game) {
	printPlayers((*g).Players)
	winners, ties, losers := g.DetermineWinners()
	for _, player := range winners {
		fmt.Printf("%s won\n", player.Name)
		g.PayoutWin(player)
	}
	for _, player := range ties {
		fmt.Printf("%s tied\n", player.Name)
		g.TiedDealer(player)
	}
	for _, player := range losers {
		fmt.Printf("%s lost\n", player.Name)
		g.TakeBet(player)
	}
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

func acceptBet(g *blackjack.Game, p *blackjack.Player) {
	noBetPlaced := true
	for noBetPlaced {
		fmt.Print("Enter bet amount: ")
		reader := bufio.NewReader(os.Stdin)
		betStr, _ := reader.ReadString('\n')
		betStr = strings.Replace(betStr, "\n", "", -1)
		bet, _ := strconv.Atoi(betStr)
		noBetPlaced = !g.PlaceBet(p, bet)
		if noBetPlaced {
			fmt.Println("Invalid Bet")
		}
	}
}

func printPlayers(players []blackjack.Player) {
	for _, player := range players {
		fmt.Println(player)
	}
}
