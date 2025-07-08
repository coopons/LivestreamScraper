package model

import (
	"fmt"
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

func (k KickStream) ToModelStream() Stream {
	return Stream{
		ID:           fmt.Sprintf("%d", k.ChannelID),
		UserID:       fmt.Sprintf("%d", k.BroadcasterUserID),
		UserName:     k.Slug,
		Title:        k.StreamTitle,
		GameID:       fmt.Sprintf("%d", k.Category.ID),
		GameName:     k.Category.Name,
		Language:     k.Language,
		ViewerCount:  k.ViewerCount,
		StartedAt:    k.StartedAt, // Already in RFC3339 format
		ThumbnailURL: k.Thumbnail,
		IsMature:     k.HasMatureContent,
	}
}