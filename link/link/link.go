package link

// Accepts html as an io.Reader

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"strings"
)

// Test this is a test
func Test() {
	fmt.Println("Test")
}

// Link A class representing a link and the text associated with it in html
type Link struct {
	Href string
	Text string
}

// ParseLinks Parse Links for some html
func ParseLinks(htmlText io.Reader) []Link {
	var links []Link
	htmlNode, err := html.Parse(htmlText)
	if err != nil {
		panic(err)
	}

	dfs(htmlNode, "", false, &links)
	return links
}

func dfs(node *html.Node, padding string, underHref bool, links *[]Link) string {
	isHrefNode := false
	hrefLink := ""
	hrefText := ""
	if node.Type == html.TextNode {
		hrefText = strings.TrimSpace(node.Data) + " "
	}
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, s := range node.Attr {
			if s.Key == "href" {
				isHrefNode = true
				underHref = true
				hrefLink = s.Val
				// links = append(links, Link{s.Val, s.Val})
			}
		}
	}
	// dfs the html tree
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		hrefText = hrefText + "" + dfs(c, padding+"  ", underHref, links)
	}
	if isHrefNode {
		// Append and assign to the value that links points at
		*links = append(*links, Link{hrefLink, strings.TrimSpace(hrefText)})
	}
	return hrefText
}
