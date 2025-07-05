package db

import (
	"context"
	"log"

	"github.com/coopons/livestream_scraper/internal/model"
)

// Saves the stream information in the DB
func SaveStream(s model.Stream, platform string) error {
	exists, err := StreamExists(platform, s.ID)
	if err != nil {
		return err
	}
	if exists {
		return nil // Stream already in db
	}

	_, err = Pool.Exec(context.Background(),
		`INSERT INTO streams (platform, stream_id, streamer_name, title, game, language, thumbnail_url, is_mature, started_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT DO NOTHING`,
		platform, s.ID, s.UserName, s.Title, s.GameName, s.Language, s.ThumbnailURL, s.IsMature, s.StartedAt,
	)

	if err != nil {
		log.Println("Failed to insert stream:", err)
	}
	return err
}

// Saves the relevant stream information as a snapshot
func SaveSnapshot(s model.Stream) error {
	_, err := Pool.Exec(context.Background(),
	`INSERT INTO stream_snapshots (stream_id, viewer_count, is_live)
	VALUES ($1, $2, true)`,
	s.ID, s.ViewerCount,
	)

	if err != nil {
		log.Println("Failed to insert snapshot:", err)
	}
	return err
}