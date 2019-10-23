package storyhtml

import (
	"fmt"
	"github.com/gtcooke94/gophercises/cyoa/cyoa"
	"html/template"
	"net/http"
	"strconv"
	"strings"
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
		<div><a href="-1"><p padding-top: 3cm;>Restart</p></div>
	</body>
</html>`

var tmpl *template.Template = template.Must(template.New("story_display").Parse(tmpl_html))

var option_channel chan int = make(chan int, 1)
var render_html chan bool = make(chan bool, 1)

type StoryHTML struct {
	Port      int
	h         handler
	start_arc *cyoa.StoryArc
}

func (s *StoryHTML) GetOption() int {
	option := <-option_channel
	return option
}

func (s *StoryHTML) DisplayArc(arc *cyoa.StoryArc) bool {
	// Update the arc so the other thread can handle it
	s.h.arc = arc
	render_html <- true
	return true
}

func (s *StoryHTML) StartStory(filename string) *cyoa.StoryArc {
	intro := cyoa.StartStory(filename)
	s.h = handler{intro}
	go http.ListenAndServe(fmt.Sprintf(":%d", s.Port), &s.h)
	return cyoa.StartStory(filename)
}

type handler struct {
	arc *cyoa.StoryArc
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	option_str := strings.TrimSpace(r.URL.Path)
	if option_str == "" || option_str == "/" {
		<-render_html
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
	// Send the selection option, then wait for the arc to be updated in the main thread so that the html can be rendered
	option_channel <- option
	<-render_html
	err = tmpl.Execute(w, h.arc)
	if err != nil {
		panic(err)
	}
}
