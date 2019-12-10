// Simple HTTP server with regex-based router

package main

import (
	"fmt"
	"net/http"
	"regexp"
)

type Route struct {
	Regex   *regexp.Regexp
	Handler func(w http.ResponseWriter, r *http.Request)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "HOME")
}

func hotelHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "HOTEL %s", r.URL.Path)
}

func main() {
	routes := []Route{
		{regexp.MustCompile(`^/$`), homeHandler},
		{regexp.MustCompile(`^/hotels/(\d+)$`), hotelHandler},
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for _, route := range routes {
			matches := route.Regex.FindStringSubmatch(r.URL.Path)
			if len(matches) >= 1 {
				route.Handler(w, r)
				return
			}
		}
		http.NotFound(w, r)
	})
	fmt.Println("listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
