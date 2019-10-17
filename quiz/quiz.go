package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

var mux sync.Mutex = sync.Mutex{}

func main() {
	timer_length_ptr := flag.Int("time", 5, "Time length of quiz")
	flag.Parse()
	file, _ := os.Open("problems.csv")
	r := csv.NewReader(file)
	lines, _ := r.ReadAll()
	num_correct := 0
	num_asked := 0
	total := len(lines)
	fmt.Println("Press enter when you are ready for the quiz to start")
	start_reader := bufio.NewReader(os.Stdin)
	start_reader.ReadString('\n')
	ch := make(chan int)
	timerchan := make(chan bool)
	go run_quiz(lines, ch)
	go run_timer(timerchan, *timer_length_ptr)
bigloop:
	for {
		select {
		case <-timerchan:
			break bigloop
		case correct := <-ch:
			num_correct += correct
			num_asked++
			if num_asked == total {
				break bigloop
			}
		}
	}
	fmt.Printf("%v/%v Correct", num_correct, total)
	return
}

func run_timer(timechan chan bool, timer_length int) {
	time.Sleep(time.Duration(timer_length) * time.Second)
	timechan <- true
}

func run_quiz(lines [][]string, ch chan int) {
	for _, line := range lines {
		question := line[0]
		answer := line[1]
		fmt.Println(question)
		answer_reader := bufio.NewReader(os.Stdin)
		input, _ := answer_reader.ReadString('\n')
		input = strings.TrimSuffix(input, "\n")
		if answer == input {
			ch <- 1
		} else {
			ch <- 0
		}
	}
}
