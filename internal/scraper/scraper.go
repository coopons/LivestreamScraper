package scraper

import "github.com/coopons/livestream_scraper/internal/model"

type StreamCollector interface {
	// Fetch live streams up to limit
	GetLiveStreams(limit int) ([]model.Stream, error)
	Platform() string
}