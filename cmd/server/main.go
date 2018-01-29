package main

import (
	"net/http"
	"log"
	"flag"

	"github.com/lucasponce/swiftsunshine/prometheus"
)

func main() {
	bindAddr := flag.String("bindAddr", ":8080", "Address to bind to for serving")
	prometheusAddr := flag.String("prometheusAddr", "http://localhost:9090", "Address of prometheus instance for graph generation")
	flag.Parse()

	http.Handle("/", prometheus.NewQueryHandler(*prometheusAddr))
	log.Fatal(http.ListenAndServe(*bindAddr, nil))
}
