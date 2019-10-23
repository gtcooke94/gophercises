package cyoa

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func Print() {
	fmt.Println("in cyoa")
}

type StoryArc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}

type Option struct {
	Text      string `json:"text"`
	ArcString string `json:"arc"`
	Arc       *StoryArc
}

// Keep a global variable that is the story's start
var start_arc_ptr *StoryArc

func StartStory(filename string) *StoryArc {
	story_json := parse_json(filename)
	story_start := build_story(story_json)
	// fmt.Println(story_start)
	return story_start
}

func (arc *StoryArc) Advance(selection int) *StoryArc {
	// If the selection is -1, return to the beginning
	if selection == -1 {
		return start_arc_ptr
	}
	return arc.Options[selection].Arc
}

type StoryJson map[string]*StoryArc

func build_story(story_json *StoryJson) *StoryArc {
	start_arc := (*story_json)["intro"]
	start_arc_ptr = start_arc
	// Link the arcs
	for k, _ := range *story_json {
		for i, _ := range (*story_json)[k].Options {
			option := &((*story_json)[k].Options[i])
			link_arc(story_json, i, option)
		}
	}
	return start_arc
}

func link_arc(story_json *StoryJson, i int, option *Option) {
	option.Arc = (*story_json)[option.ArcString]
}

func parse_json(filename string) *StoryJson {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	bytes, _ := ioutil.ReadAll(file)
	var story_json StoryJson
	json.Unmarshal(bytes, &story_json)
	return &story_json
}

type StoryRunner interface {
	GetOption() int
	DisplayArc(*StoryArc) bool
	StartStory(string) *StoryArc
}
