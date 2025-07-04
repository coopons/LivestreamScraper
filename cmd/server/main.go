package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/coopons/livestream_scraper/internal/api"
	"github.com/coopons/livestream_scraper/internal/db"
	"github.com/coopons/livestream_scraper/internal/web"
)

func main() {
	mux := http.NewServeMux()
	db.InitDB()

	clientID := os.Getenv("TWITCH_CLIENT_ID")
	clientSecret := os.Getenv("TWITCH_CLIENT_SECRET")

	mux.HandleFunc("/", web.HomeHandler)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("internal/web/static"))))
	mux.HandleFunc("/api/snapshots", web.SnapshotDataHandler)

	go func() {
		ticker := time.NewTicker(10 * time.Minute)	// Runs the Collector every 5 minutes
		defer ticker.Stop()
		
		for {
			err := runCollection(clientID, clientSecret)
			if err != nil {
				log.Println("Collection error:", err)
			}
			<-ticker.C
		}
	}()
	
	
	log.Println("Server running on http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func runCollection(clientID, clientSecret string) error {
	token, err := api.GetAppAccessToken(clientID, clientSecret)
	if err != nil {
		return err
	}

	streams, err := api.GetAllLiveStreams(clientID, token, 200) // Only get top 200 streams to start out
	if err != nil {
		return err
	}

	for _, s := range streams {
		if err := db.SaveStream(s, "twitch"); err != nil {
			log.Println("SaveStream error:", err)
		}
		if err := db.SaveSnapshot(s); err != nil {
			log.Println("SaveSnapshot error:", err)
		}
	}
	return nil
}