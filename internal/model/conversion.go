package model

import (
	"time"
)

func (yt YtStream) ToModelStream() Stream {
	startedAt := time.Unix(yt.StartTime, 0).Format(time.RFC3339)
	gameName := ""
	if len(yt.Categories) > 0 {
		gameName = yt.Categories[0]
	}

	return Stream{
		ID:           yt.ID,
		UserID:       yt.UserID,
		UserName:     yt.UserName,
		Title:        yt.Title,
		GameID:       "", // YouTube doesn't provide a numeric game ID
		GameName:     gameName,
		Language:     yt.Language,
		ViewerCount:  yt.ViewCount,
		StartedAt:    startedAt,
		ThumbnailURL: yt.Thumbnail,
		IsMature:     false, // YouTube doesn't expose this
	}
}