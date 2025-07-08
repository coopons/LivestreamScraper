package scraper

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"

	"github.com/coopons/livestream_scraper/internal/model"
)

type YoutubeScraper struct {}

func (y *YoutubeScraper) GetLiveStreams(limit int) ([]model.Stream, error) {
	return ScrapeYoutubeLivestreams()
}

func (y *YoutubeScraper) Platform() string {
	return "youtube"
}

func ScrapeYoutubeLivestreams() ([]model.Stream, error) {
	cmd := exec.Command("yt-dlp", "--dump-json", "https://www.youtube.com/results?search_query=live&sp=EgJAAQ%253D%253D")
	
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start yt-dlp command: %w", err)
	}

	var allStreams []model.Stream
	decoder := json.NewDecoder(stdout)

	for decoder.More() {
		var ytStream model.YtStream
		if err := decoder.Decode(&ytStream); err != nil {
			fmt.Println("Error decoding youtube Lives JSON:", err)
			continue
		}

		if !ytStream.LiveStatus {
			log.Printf("Skipping non-live video: %s (%s)\n", ytStream.Title, ytStream.ID)
			continue
		}

		allStreams = append(allStreams, ytStream.ToModelStream())
	}

	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("yt-dlp command failed: %w", err)
	}
	
	return allStreams, nil
}
