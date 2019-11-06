package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/gtcooke94/gophercises/link/link"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	site := flag.String("s", "https://gophercises.com", "The website you want to build a sitemap for")
	maxDepth := flag.Int("d", -1, "Max Depth in website to explore")
	flag.Parse()
	l := link.Link{Href: *site, Text: ""}
	links := crawlSite(l, *site, *maxDepth)
	linksToXML(links)
}

func linksToXML(links map[string]int) {
	f, err := os.Create("output.xml")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var urls urlset
	for url := range links {
		urls.Urls = append(urls.Urls, loc{url})
	}
	urls.Xmlns = xmlns

	fmt.Fprint(f, xml.Header)
	enc := xml.NewEncoder(f)
	enc.Indent("", "  ")
	enc.Encode(urls)
}

type linkQueueElement struct {
	l     link.Link
	depth int
}

const xmlns string = "http://www.sitemaps.org/schemas/sitemap/0.9"

type urlset struct {
	Xmlns string `xml:"xmlns,attr"`
	Urls  []loc  `xml:"url"`
}

type loc struct {
	Address string `xml:"loc"`
}

func crawlSite(l link.Link, site string, maxDepth int) map[string]int {
	var linksQueue []linkQueueElement
	var curLink link.Link
	linksQueue = append(linksQueue, linkQueueElement{l, 0})
	visitedLinks := map[string]int{}
	curDepth := 0
	for len(linksQueue) != 0 {
		curDepth = linksQueue[0].depth
		curLink = linksQueue[0].l
		if maxDepth != -1 && curDepth > maxDepth {
			return visitedLinks
		}
		linksQueue = linksQueue[1:]
		// Handle current link
		if _, ok := visitedLinks[curLink.Href]; ok {
			// Continue the loop, we've already visited this link
			continue
		}
		newLinks := link.ParseLinks(getBody(curLink))
		visitedLinks[curLink.Href] = curDepth
		// Add it's links to queue
		for _, l := range newLinks {
			formattedLink, ok := buildLink(l, site)
			if ok {
				linksQueue = append(linksQueue, linkQueueElement{formattedLink, curDepth + 1})
			}
		}
	}
	return visitedLinks
}

func buildLink(l link.Link, site string) (link.Link, bool) {
	newLink := l
	returnFlag := false
	if strings.HasPrefix(l.Href, site) {
		returnFlag = true
	} else if strings.HasPrefix(l.Href, "/") {
		newLink = link.Link{Href: site + l.Href}
		returnFlag = true
	}
	curHref := newLink.Href
	if strings.HasSuffix(curHref, "/") {
		newLink = link.Link{Href: curHref[:len(curHref)-1], Text: ""}
	}
	return newLink, returnFlag
}

func getBody(l link.Link) io.Reader {
	resp, err := http.Get(l.Href)
	if err != nil {
		panic(err)
	}
	return resp.Body
}
