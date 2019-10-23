package storyhtml

import (
	"fmt"
	"github.com/gtcooke94/gophercises/cyoa/cyoa"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	// "os"
	// "time"
)

func main() {
	fmt.Println("vim-go")
}

const tmpl_html = `
<!DOCTYPE html>
<html>
	<head>
		<link rel="icon" href="data:,">
		<meta charset="utf-8">	
		<title>Choose Your Own Adventure</title>
	</head>
	<body>
		<h1>{{ .Title }}</h1>
		{{range .Story}}<p>{{ . }}</p>{{end}}
	</body>
	<body>
		{{ range $i, $val := .Options }}<div><a href="{{ $i }}">{{ $i }}: {{ $val.Text }}</div>{{end}}
	</body>
	<body>
		<div><a href="-1">Restart</div>
	</body>
</html>`

var tmpl *template.Template = template.Must(template.New("story_display").Parse(tmpl_html))

// TODO use this channel in combo with the ServeHTTP, ListenAndServe stuff to get it working with the layout I have in the main runner file
var option_channel chan int = make(chan int, 1)
var render_html chan bool = make(chan bool, 1)
var finished chan bool = make(chan bool)

type StoryHTML struct {
	Port      int
	h         handler
	start_arc *cyoa.StoryArc
}

func (s *StoryHTML) GetOption() int {
	fmt.Println("")
	option := <-option_channel
	fmt.Println("STORY option RECV")
	return option
}

func (s *StoryHTML) DisplayArc(arc *cyoa.StoryArc) bool {
	s.h.arc = arc
	render_html <- true
	fmt.Println("STORY html SEND")
	return true
}

func (s *StoryHTML) StartStory(filename string) *cyoa.StoryArc {
	intro := cyoa.StartStory(filename)
	s.h = handler{intro}
	// option_channel <-
	go http.ListenAndServe(fmt.Sprintf(":%d", s.Port), &s.h)
	return cyoa.StartStory(filename)
}

type handler struct {
	arc *cyoa.StoryArc
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// This gets the next option and also serves the current option
	option_str := strings.TrimSpace(r.URL.Path)
	// fmt.Printf("option_str is %v, length %v\n", option_str, len(option_str))
	if option_str == "" || option_str == "/" {
		<-render_html
		fmt.Println("SERVE render_html RECV")
		err := tmpl.Execute(w, h.arc)
		if err != nil {
			panic(err)
		}
		return
	}
	option_str = option_str[1:]
	option, err := strconv.Atoi(option_str)
	if err != nil {
		panic(err)
	}
	option_channel <- option
	fmt.Println("SERVE option SEND")
	<-render_html
	fmt.Println("SERVE render_html RECV")
	err = tmpl.Execute(w, h.arc)
	if err != nil {
		panic(err)
	}
}

// type StoryRunner interface {
//     GetOption() int
//     DisplayArc(*StoryArc) bool
//     StartStory(string) *StoryArc
// }
