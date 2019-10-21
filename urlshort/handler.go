package urlshort

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"net/http"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	//	TODO: Implement this...
	// http.HandlerFunc is a type that is really just a function
	return func(w http.ResponseWriter, r *http.Request) {
		asked_path := r.URL.Path
		new_path, ok := pathsToUrls[asked_path]
		if ok {
			fmt.Printf("Redirecting %v to %v\n", asked_path, new_path)
			// 302 is the integer code for Found in redirect
			// Better to use http.StatusFound
			// http.Redirect(w, r, new_path, 302)
			http.Redirect(w, r, new_path, http.StatusFound)
			return
		}
		fmt.Printf("Serving %v\n", asked_path)
		fallback.ServeHTTP(w, r)
	}
	// for shortened, redirect := range pathsToUrls {
	//     fmt.Printf("%v, %v", shortened, redirect)
	// }
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// TODO: Implement this...
	// Marshal take data into yaml
	// Unmarsahl take data out of yaml

	// Read yaml
	var pathUrls []ymlUrls
	err := yaml.Unmarshal(yml, &pathUrls)
	if err != nil {
		return nil, err
	}

	// Make a map
	pathsToUrls := map[string]string{}
	for _, pathurl := range pathUrls {
		pathsToUrls[pathurl.Path] = pathurl.URL
		fmt.Printf("%v, %v\n", pathurl.Path, pathurl.URL)
	}

	// Return the MapHandler
	return MapHandler(pathsToUrls, fallback), nil
}

type ymlUrls struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}
