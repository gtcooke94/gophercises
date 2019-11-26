package main

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/gtcooke94/gophercises/quiet_hn/hn"
)

func main() {
	// parse flags
	var port, numStories int
	flag.IntVar(&port, "port", 3000, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.Parse()

	tpl := template.Must(template.ParseFiles("./index.gohtml"))

	http.HandleFunc("/", handler(numStories, tpl))

	// Start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handler(numStories int, tpl *template.Template) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		stories, err := getTopItems(numStories)
		data := templateData{
			Stories: stories,
			Time:    time.Now().Sub(start),
		}
		err = tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
}

func getTopItems(numStories int) ([]item, error) {
	var client hn.Client
	ids, err := client.TopItems()
	stories := make(chan concurrencyItem, len(ids))
	if err != nil {
		return nil, errors.New("Failed to load top stories")
	}

	counter := make(chan int, len(ids))
	for i, id := range ids {
		go func(id int, order int) {
			hnItem, err := client.GetItem(id)
			if err != nil {
				counter <- 0
				return
			}
			item := parseHNItem(hnItem)
			if isStoryLink(item) {
				stories <- concurrencyItem{item: item, orderBy: order}
			}
			counter <- 0
			return
		}(id, i)
	}
	for len(counter) != len(ids) {
		time.Sleep(1 * time.Millisecond)
	}
	// At this point all stories are processed
	close(stories)
	// Put them all in a list
	storiesList := make([]concurrencyItem, 0, len(ids))
	for story := range stories {
		storiesList = append(storiesList, story)
	}
	// Sort the stories by orderBy which should keep them in order
	sort.Slice(storiesList, func(i, j int) bool {
		return storiesList[i].orderBy < storiesList[j].orderBy
	})
	orderedStories := make([]item, 0, numStories)
	// Get then return the top 30
	for _, story := range storiesList[:30] {
		orderedStories = append(orderedStories, story.item)
	}

	return orderedStories, nil
}

func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(hnItem hn.Item) item {
	ret := item{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}

// item is the same as the hn.Item, but adds the Host field
type item struct {
	hn.Item
	Host string
}

type concurrencyItem struct {
	item
	orderBy int
}

type templateData struct {
	Stories []item
	Time    time.Duration
}
