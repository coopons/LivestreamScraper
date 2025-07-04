package web

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/coopons/livestream_scraper/internal/api"
	"github.com/coopons/livestream_scraper/internal/db"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	clientID := os.Getenv("TWITCH_CLIENT_ID")
	clientSecret := os.Getenv("TWITCH_CLIENT_SECRET")

	token, err := api.GetAppAccessToken(clientID, clientSecret)
	if err != nil {
		http.Error(w, "Failed to authenticate", http.StatusInternalServerError)
		return
	}

	streams, err := api.GetLiveStreams(clientID, token, 12)
	if err != nil {
		http.Error(w, "Failed to get live streams", http.StatusInternalServerError)
	}
	
	tmpl := template.Must(template.ParseFiles("internal/web/templates/index.html"))
	err = tmpl.Execute(w, streams)
	if err != nil {
		log.Println("Template error:", err)
	}
}

func SnapshotDataHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Pool.Query(context.Background(),
		`SELECT stream_id, viewer_count, timestamp FROM stream_snapshots
		ORDER BY timestamp ASC LIMIT 1000`)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Point struct {
		StreamID 	string	 	`json:"stream_id"`
		ViewerCount	int			`json:"viewer_count"`
		Timestamp	time.Time	`json:"timestamp"`
	}
	var points []Point

	for rows.Next() {
		var p Point
		if err := rows.Scan(&p.StreamID, &p.ViewerCount, &p.Timestamp); err == nil {
			points = append(points, p)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(points)
}