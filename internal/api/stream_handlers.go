package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/coopons/livestream_scraper/internal/db"
	"github.com/coopons/livestream_scraper/internal/model"
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

func GetTopRecentStreams(limit int) ([]model.Stream, error) {
	rows, err := db.Pool.Query(context.Background(),
		`SELECT *
		FROM (
			SELECT DISTINCT ON (streamer_name)
				platform,
				stream_id,
				streamer_name,
				title,
				game,
				language,
				thumbnail_url,
				timestamp,
				viewer_count
			FROM (
				SELECT DISTINCT ON (s.stream_id)
					s.platform,
					s.stream_id,
					s.streamer_name,
					s.title,
					s.game,
					s.language,
					s.thumbnail_url,
					ss.timestamp,
					ss.viewer_count
				FROM
					streams s
				JOIN
					stream_snapshots ss ON s.stream_id = ss.stream_id
				WHERE
					ss.timestamp > NOW() - INTERVAL '10 minutes'
				ORDER BY
					s.stream_id,
					ss.timestamp DESC
			) AS latest_streams
			ORDER BY
				streamer_name,
				timestamp DESC
		) AS latest_per_streamer
		ORDER BY
			viewer_count DESC
		LIMIT $1;`, limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var streams []model.Stream
	var timestamp time.Time
	for rows.Next() {
		var s model.Stream
		err := rows.Scan(
				&s.Platform,
				&s.ID,
				&s.UserName,
				&s.Title,
				&s.GameName,
				&s.Language,
				&s.ThumbnailURL,
				&timestamp,
				&s.ViewerCount)
		if err != nil {
			return nil, err
		}
		streams = append(streams, s)
	}

	return streams, nil
}