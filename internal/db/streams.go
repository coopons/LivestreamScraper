package db

import (
	"context"
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