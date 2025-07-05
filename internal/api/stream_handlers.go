package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/coopons/livestream_scraper/internal/db"
)

func SnapshotDataHandler(w http.ResponseWriter, r *http.Request) {
	streamID := r.URL.Query().Get("stream_id")
	if streamID == "" {
		http.Error(w, "Missing stream_id", http.StatusBadRequest)
		return
	}
	
	// Get the streamer from current streamID
	var streamerName string
	err := db.Pool.QueryRow(context.Background(),
		`SELECT streamer_name 
		FROM streams 
		WHERE stream_id = $1`, streamID).Scan(&streamerName)
	if err != nil {
		http.Error(w, "Could not find streamer", http.StatusInternalServerError)
		return
	}

	// Get 5 most recent streams
	rows, err := db.Pool.Query(context.Background(),
		`SELECT stream_id
		FROM streams
		WHERE streamer_name = $1
		ORDER BY started_at DESC
		LIMIT 5`, streamerName)
	if err != nil {
		log.Println("Could not locate recent streams.")
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var recentStreamIDs []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err == nil {
			recentStreamIDs = append(recentStreamIDs, id)
		}
	}

	type Point struct {
		StreamID 	string	 	`json:"stream_id"`
		ViewerCount	int			`json:"viewer_count"`
		Timestamp	time.Time	`json:"timestamp"`
	}
	
	type StreamSnapshots struct {
		StreamID 		string  `json:"stream_id"`
		Snapshots	   []Point	`json:"snapshots"`
		AverageViewers 	int		`json:"average_viewers"`
		DurationMinutes	int		`json:"duration_minutes"`
	}
	
	var response []StreamSnapshots
	// Get snapshots for the 5 recent streams
	for _, sid := range recentStreamIDs {
		snapRows, err := db.Pool.Query(context.Background(),
			`SELECT stream_id, viewer_count, timestamp
			FROM stream_snapshots
			WHERE stream_id = $1
			ORDER BY timestamp ASC
			LIMIT 1000`, sid)
		if err != nil {
			log.Println("Error in retrieving snapshots.")
			http.Error(w, "DB Error", http.StatusInternalServerError)
			return
		}
		
		var snaps []Point
		for snapRows.Next() {
			var p Point
			if err := snapRows.Scan(&p.StreamID, &p.ViewerCount, &p.Timestamp); err == nil {
				snaps = append(snaps, p)
			}
		}
		snapRows.Close()

		// Calculate average viewers and stream duration
		var totalViewers int
		var startTime, endTime time.Time
		for i, p := range snaps {
			totalViewers += p.ViewerCount
			if i == 0 {
				startTime = p.Timestamp
			}
			endTime = p.Timestamp
		}

		duration := endTime.Sub(startTime)
		average := 0
		if len(snaps) > 0 {
			average = totalViewers / len(snaps)
		}
		
		response = append(response, StreamSnapshots{
			StreamID: sid,
			Snapshots: snaps,
			AverageViewers: average,
			DurationMinutes: int(duration),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
