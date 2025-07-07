package db

import (
	"context"
	"time"
)

func StreamExists(platform, streamID string) (bool, error) {
	var exists bool
	err := Pool.QueryRow(
		context.Background(),
		`SELECT EXISTS (
			SELECT 1 FROM streams WHERE platform = $1 AND stream_id = $2
		)`,
		platform, streamID,
	).Scan(&exists)

	return exists, err
}

func GetLatestSnapshotTime() (time.Time, error) {
	var latest time.Time
	err := Pool.QueryRow(context.Background(),
		`SELECT MAX(timestamp) FROM stream_snapshots`).Scan(&latest)
	if err != nil {
		return time.Time{}, err
	}
	return latest, nil
}