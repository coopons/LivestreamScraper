package web

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/coopons/livestream_scraper/internal/scraper"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	clientID := os.Getenv("TWITCH_CLIENT_ID")
	clientSecret := os.Getenv("TWITCH_CLIENT_SECRET")

	token, err := scraper.GetAppAccessToken(clientID, clientSecret)
	if err != nil {
		http.Error(w, "Failed to authenticate", http.StatusInternalServerError)
		return
	}

	streams, err := scraper.GetLiveStreams(clientID, token, 12)
	if err != nil {
		http.Error(w, "Failed to get live streams", http.StatusInternalServerError)
	}
	
	tmpl := template.Must(template.ParseFiles("internal/web/templates/index.html"))
	err = tmpl.Execute(w, streams)
	if err != nil {
		log.Println("Template error:", err)
	}
}