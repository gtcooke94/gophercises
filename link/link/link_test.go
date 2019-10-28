package link

import (
	"strings"
	"testing"
)

func TestParseLinks(t *testing.T) {
	outputs := make([][]Link, len(htmls))
	for i, html := range htmls {
		outputs[i] = ParseLinks(strings.NewReader(html))
	}
	for i, _ := range solns {
		for j, _ := range solns[i] {
			if outputs[i][j] != solns[i][j] {
				t.Errorf("Example %v Failed. Expected %v, got %v\n", i, solns[i][j], outputs[i][j])
			}
		}
	}
}

var html1 = `
<html>
<body>
  <h1>Hello!</h1>
  <a href="/other-page">A link to another page</a>
</body>
</html>`

var html2 = `
<html>
<head>
  <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/font-awesome/4.7.0/css/font-awesome.min.css">
</head>
<body>
  <h1>Social stuffs</h1>
  <div>
    <a href="https://www.twitter.com/joncalhoun">
      Check me out on twitter
      <i class="fa fa-twitter" aria-hidden="true"></i>
    </a>
    <a href="https://github.com/gophercises">
      Gophercises is on <strong>Github</strong>!
    </a>
  </div>
</body>
</html>
`

var htmls = []string{html1, html2}

var soln1 = []Link{Link{"/other-page", "A link to another page"}}
var soln2 = []Link{
	Link{"https://www.twitter.com/joncalhoun", "Check me out on twitter"},
	Link{"https://github.com/gophercises", "Gophercises is on Github !"},
}
var solns = [][]Link{soln1, soln2}
