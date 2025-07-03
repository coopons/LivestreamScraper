package main

import (
	"log"
	"net/http"

	"github.com/coopons/livestream_scraper/internal/web"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", web.HomeHandler)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("internal/web/static"))))

	log.Println("Server running on https://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}