package main

import (
	"flag"
	"fmt"
	"github.com/gtcooke94/gophercises/cyoa/cyoa"
	"github.com/gtcooke94/gophercises/cyoa/story_cli"
)

func main() {
	serve_method := flag.String("s", "cli", "Method you want the story served via")
	flag.Parse()
	fmt.Println("Start")
	// story_json := cyoa.StartStory("gopher.json")
	// story_start := cyoa.StartStory("gopher.json")
	var story_runner cyoa.StoryRunner
	if *serve_method == "cli" {
		story_runner = storycli.StoryCLI{}
	}
	story_start := story_runner.StartStory("gopher.json")
	current_arc := story_start
	for {
		end_flag := story_runner.DisplayArc(current_arc)
		if !end_flag {
			break
		}
		selection := story_runner.GetOption()
		current_arc = current_arc.Advance(int(selection))
	}
}
