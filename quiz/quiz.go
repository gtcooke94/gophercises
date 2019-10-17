package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

var mux sync.Mutex = sync.Mutex{}

func main() {
	fmt.Println("Go!")
	file, _ := os.Open("problems.csv")
	r := csv.NewReader(file)
	lines, _ := r.ReadAll()
	num_correct := 0
	total := len(lines)
	fmt.Println("Press enter when you are ready for the quiz to start")
	start_reader := bufio.NewReader(os.Stdin)
	start_reader.ReadString('\n')
	ch := make(chan int)
	timerchan := make(chan bool)
	go run_quiz(lines, ch)
	go run_timer(timerchan)
bigloop:
	for {
		select {
		case <-timerchan:
			break bigloop
		case num_correct = <-ch:
			if num_correct == total {
				break bigloop
			}
		}
	}
	fmt.Printf("%v/%v Correct", num_correct, total)
	return
}

func run_timer(timechan chan bool) {
	time.Sleep(5 * time.Second)
	timechan <- true
}

func run_quiz(lines [][]string, ch chan int) {
	num_correct := 0
	for _, line := range lines {
		question := line[0]
		answer := line[1]
		fmt.Println(question)
		answer_reader := bufio.NewReader(os.Stdin)
		input, _ := answer_reader.ReadString('\n')
		input = strings.TrimSuffix(input, "\n")
		if answer == input {
			num_correct++
			ch <- num_correct
		}
	}
}
