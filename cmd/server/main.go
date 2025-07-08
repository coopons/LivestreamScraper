package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/coopons/livestream_scraper/internal/api"
	"github.com/coopons/livestream_scraper/internal/db"
	"github.com/coopons/livestream_scraper/internal/web"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, using environment variables directly")
	}
}

func main() {
	mux := http.NewServeMux()
	db.InitDB()

	twitchClientID := os.Getenv("TWITCH_CLIENT_ID")
	twitchClientSecret := os.Getenv("TWITCH_CLIENT_SECRET")
	kickClientID := os.Getenv("KICK_CLIENT_ID")
	kickClientSecret := os.Getenv("KICK_CLIENT_SECRET")

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("internal/web/static"))))
	mux.HandleFunc("/", web.HomeHandler)
	mux.HandleFunc("/api/snapshots", api.SnapshotDataHandler)
	mux.HandleFunc("/api/control", api.ControlHandler)
	mux.HandleFunc("/api/next-run", api.NextRunHandler)

	api.StartCollector(twitchClientID, twitchClientSecret, kickClientID, kickClientSecret, 10*time.Minute)
	
	log.Println("Server running on http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}

