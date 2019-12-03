package main

import (
	// "errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
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
	sc := storyCache{
		numStories: numStories,
		duration:   3 * time.Second,
	}

	go func() {
		ticker := time.NewTicker(3 * time.Second)
		for {
			temp := storyCache{
				numStories: numStories,
				duration:   6 * time.Second,
			}
			temp.stories()
			sc.mutex.Lock()
			sc.cache = temp.cache
			sc.expiration = temp.expiration
			sc.mutex.Unlock()
			<-ticker.C
		}
	}()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		stories, _ := sc.stories()
		data := templateData{
			Stories: stories,
			Time:    time.Now().Sub(start),
		}
		err := tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
}

type storyCache struct {
	cache      []item
	useA       bool
	expiration time.Time
	duration   time.Duration
	mutex      sync.Mutex
	numStories int
}

func (sc *storyCache) stories() ([]item, error) {
	// Use pointer to storyCache so that the mutex doesn't get copied, etc
	// Using global cache
	sc.mutex.Lock()
	// Good practice to defer unlocks, that way you can't forget if you have multiple returns
	defer sc.mutex.Unlock()
	if sc.expiration.Sub(time.Now()) > 0 {
		return sc.cache, nil
	}

	stories, err := getStories(sc.numStories)
	if err != nil {
		return nil, err
	}
	sc.cache = stories
	sc.expiration = time.Now().Add(sc.duration)
	return sc.cache, nil
}

// var (
//     cache           []item
//     cacheExpiration time.Time
//     cacheMutex      sync.Mutex
// )

// func getCachedStories(numStories int) ([]item, error) {
//     // Using global cache
//     cacheMutex.Lock()
//     // Good practice to defer unlocks, that way you can't forget if you have multiple returns
//     defer cacheMutex.Unlock()
//     if cacheExpiration.Sub(time.Now()) > 0 {
//         return cache, nil
//     }
//
//     stories, err := getStories(numStories)
//     if err != nil {
//         return nil, err
//     }
//     cache = stories
//     cacheExpiration = time.Now().Add(15 * time.Minute)
//     return stories, nil
// }

func getStories(numStories int) ([]item, error) {
	var client hn.Client
	ids, err := client.TopItems()
	stories := make([]item, 0)

	currentIdIdx := 0
	numToGet := int(numStories * 5 / 4)
	for len(stories) < numStories {
		newStories, _ := getTopItems(ids[currentIdIdx:currentIdIdx+numToGet], numStories)
		stories = append(stories, newStories...)
		currentIdIdx += numToGet
		numToGet = int((numStories - len(stories)) * 5 / 4)
	}
	stories = stories[:numStories]
	return stories, err
}

func getTopItems(ids []int, numStories int) ([]item, error) {
	stories := make(chan concurrencyItem, len(ids))
	for i, id := range ids {
		go func(id int, order int) {
			var client hn.Client
			hnItem, _ := client.GetItem(id)
			item := parseHNItem(hnItem)
			stories <- concurrencyItem{item: item, orderBy: order}
			return
		}(id, i)
	}
	// Put them all in a list
	storiesList := make([]concurrencyItem, 0, len(ids))
	for i := 0; i < len(ids); i++ {
		storiesList = append(storiesList, <-stories)
	}

	// Sort the stories by orderBy which should keep them in order
	sort.Slice(storiesList, func(i, j int) bool {
		return storiesList[i].orderBy < storiesList[j].orderBy
	})
	orderedStories := make([]item, 0, numStories)
	// Get then return the top 30
	for _, story := range storiesList {
		item := story.item
		if isStoryLink(item) {
			orderedStories = append(orderedStories, item)
		}
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
