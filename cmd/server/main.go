package main

import (
	"net/http"
	"log"
	"flag"
	"github.com/lucasponce/swiftsunshine/prometheus"
	"os"
)

type OnlyFiles struct {
	Fs http.FileSystem
}

func (fs OnlyFiles) Open(name string) (http.File, error) {
	f, err := fs.Fs.Open(name)
	if err != nil {
		return nil, err
	}
	stat, _ := f.Stat()
	if stat.IsDir() {
		return nil, os.ErrNotExist
	}
	return f, nil
}

func main() {
	bindAddr := flag.String("bindAddr", ":8080", "Address to bind to for serving")
	prometheusAddr := flag.String("prometheusAddr", "http://localhost:9090", "Address of prometheus instance for graph generation")
	webDir := flag.String("webDir", "web", "directory find assets to serve")
	flag.Parse()

	http.Handle("/", http.FileServer(&OnlyFiles{http.Dir(*webDir)}))
	http.Handle("/p8s", prometheus.NewQueryHandler(*prometheusAddr))
	log.Fatal(http.ListenAndServe(*bindAddr, nil))
}
