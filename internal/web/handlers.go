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
	streamID := r.URL.Query().Get("stream_id")
	if streamID == "" {
		http.Error(w, "Missing stream_id", http.StatusBadRequest)
		return
	}

	rows, err := db.Pool.Query(context.Background(),
		`SELECT stream_id, viewer_count, timestamp
		FROM stream_snapshots
		WHERE stream_id = $1
		ORDER BY timestamp ASC
		LIMIT 1000`, streamID)

	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var streamerName string
	err = db.Pool.QueryRow(context.Background(),
		`SELECT streamer_name FROM streams WHERE stream_id = $1`, streamID).Scan(&streamerName)
	if err != nil {
		http.Error(w, "Could not find streamer", http.StatusInternalServerError)
		return
	}

	var numStreams int
	err = db.Pool.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM streams WHERE streamer_name = $1`, streamerName).Scan(&numStreams)
	if err != nil {
		http.Error(w, "Error counting streams", http.StatusInternalServerError)
		return
	}

	type Point struct {
		StreamID 	string	 	`json:"stream_id"`
		ViewerCount	int			`json:"viewer_count"`
		Timestamp	time.Time	`json:"timestamp"`
	}

	type Response struct {
	NumStreams int     `json:"num_streams"`
	Snapshots  []Point `json:"snapshots"`
	}

	var points []Point

	for rows.Next() {
		var p Point
		if err := rows.Scan(&p.StreamID, &p.ViewerCount, &p.Timestamp); err == nil {
			points = append(points, p)
		}
	}

	response := Response{
		NumStreams: numStreams,
		Snapshots: points,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func ControlHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Query().Get("action") {
	case "stop":
		stopCollector()
	case "start":
		clientID := os.Getenv("TWITCH_CLIENT_ID")
		clientSecret := os.Getenv("TWITCH_CLIENT_SECRET")
		StartCollector(clientID, clientSecret, 10*time.Minute)
	default: 
		log.Println("Unknown action:", r.URL.Query().Get("action"))
	}
	w.WriteHeader(http.StatusOK)
}

func NextRunHandler(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		NextRun string  `json:"next_run"`
		Running	bool	`json:"running"`
	}
	
	if nextRunTime.IsZero() {
	nextRunTime = time.Now().Add(tickerInterval)
	}
	
	resp := Response{
		NextRun: nextRunTime.Format(time.RFC3339),
		Running: collectorRunning,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}