package web

import (
	"html/template"
	"log"
	"net/http"

	"github.com/coopons/livestream_scraper/internal/api"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	streams, err := api.GetTopRecentStreams(50) // Fetch the top 50 most popular streams
	if err != nil {
		http.Error(w, "Failed to load top recent streams", http.StatusInternalServerError)	
		return
	}	

	tmpl := template.Must(template.ParseFiles("internal/web/templates/index.html"))
	err = tmpl.Execute(w, streams)
	if err != nil {
		log.Println("Template error:", err)
	}
}