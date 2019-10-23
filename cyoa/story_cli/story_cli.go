package storycli

import (
	"bufio"
	"fmt"
	"github.com/gtcooke94/gophercises/cyoa/cyoa"
	"os"
	"strconv"
)

type StoryCLI struct {
}

func (*StoryCLI) StartStory(filename string) *cyoa.StoryArc {
	return cyoa.StartStory(filename)
}

func (*StoryCLI) GetOption() int {
	scanner := bufio.NewScanner(os.Stdin)
	// _ = scanner
	fmt.Printf("\nWhat option do you select? ")
	scanner.Scan()
	option := scanner.Text()
	option_int, _ := strconv.Atoi(option)
	return option_int
}

func (*StoryCLI) DisplayArc(arc *cyoa.StoryArc) bool {
	fmt.Printf("Chapter: %v\n", arc.Title)
	for _, paragraph := range arc.Story {
		fmt.Printf("\t%v\n", paragraph)
	}
	if len(arc.Options) == 0 {
		return false
	}
	fmt.Println("\nChoose an option below...\n")
	for i, option := range arc.Options {
		fmt.Printf("%v: %v\n", i, option.Text)
	}
	return true
}
