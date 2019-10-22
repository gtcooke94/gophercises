package main

import (
	"fmt"
	"github.com/gtcooke94/gophercises/cyoa/cyoa"
)

func main() {
	fmt.Println("Start")
	// story_json := cyoa.StartStory("gopher.json")
	story_start := cyoa.StartStory("gopher.json")
	fmt.Println(story_start)
}
