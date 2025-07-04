package main

import (
	"log"
	"net/http"
	"os"

	"github.com/coopons/livestream_scraper/internal/db"
	"github.com/coopons/livestream_scraper/internal/web"
)

func main() {
	mux := http.NewServeMux()
	db.InitDB()

	clientID := os.Getenv("TWITCH_CLIENT_ID")
	clientSecret := os.Getenv("TWITCH_CLIENT_SECRET")

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("internal/web/static"))))
	mux.HandleFunc("/", web.HomeHandler)
	mux.HandleFunc("/api/snapshots", web.SnapshotDataHandler)
	mux.HandleFunc("/api/control", web.ControlHandler)
	mux.HandleFunc("/api/next-run", web.NextRunHandler)

	web.StartCollector(clientID, clientSecret, 10)
	
	log.Println("Server running on http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}

