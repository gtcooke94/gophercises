package main

import (
	"flag"
	"fmt"
	"github.com/gtcooke94/gophercises/link/link"
	"os"
)

func main() {
	file := flag.String("f", "ex1.html", "The html file you want to parse")
	flag.Parse()
	html, err := os.Open(*file)
	if err != nil {
		panic(err)
	}
	links := link.ParseLinks(html)
	for _, link := range links {
		fmt.Println(link)
	}
}
