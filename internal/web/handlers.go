package web

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/coopons/livestream_scraper/internal/api"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		log.Printf("Ignoring request to: %s\n", r.URL.Path)
		return
	}
	
	streams, err := api.GetTopRecentStreams(50) // Fetch the top 50 most popular streams
	if err != nil {
		log.Printf("Error fetching top recent streams: %v\n", err)
		http.Error(w, "Failed to load top recent streams", http.StatusInternalServerError)	
		return
	}	

	tmpl := template.Must(template.ParseFiles("internal/web/templates/index.html"))
	err = tmpl.Execute(w, streams)
	if err != nil {
		log.Println("Template error:", err)
	}
}

func StatsPageHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := api.GetStatsPageData()
	if err != nil {
		http.Error(w, "Failed to fetch stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func StatsHandler(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles("internal/web/templates/stats.html"))
    err := tmpl.Execute(w, nil) // no dynamic data for now, just serve static HTML
    if err != nil {
        log.Println("Template error:", err)
        http.Error(w, "Failed to render stats page", http.StatusInternalServerError)
    }
}