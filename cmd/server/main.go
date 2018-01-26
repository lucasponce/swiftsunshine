package main

import (
	"net/http"
	"log"
	"fmt"
	"html"
	"github.com/lucasponce/swiftsunshine/version"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello Sunshine, %q \n", html.EscapeString(r.URL.Path))
		fmt.Fprintf(w, "Version [%q]", version.String())
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
