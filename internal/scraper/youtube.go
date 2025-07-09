package scraper

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	cmd := exec.Command("yt-dlp",
		"--dump-json",
		"--cookies", "cookies.txt",
		"--sleep-interval", "2",
		"--max-sleep-interval", "5",
		"https://www.youtube.com/results?search_query=live&sp=EgJAAQ%253D%253D")
	
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start yt-dlp command: %w | stderr: %s", err, stderr.String())
	}

	var allStreams []model.Stream
	decoder := json.NewDecoder(stdout)

	for decoder.More() {
		var ytStream model.YtStream
		if err := decoder.Decode(&ytStream); err != nil {
			fmt.Println("Error decoding youtube Lives JSON:", err)
			continue
		}

		// Skip if the stream is not live
		if !ytStream.LiveStatus {
			// log.Printf("Skipping non-live video: %s (%s)\n", ytStream.Title, ytStream.ID)
			continue
		}

		allStreams = append(allStreams, ytStream.ToModelStream())
	}

	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("yt-dlp command failed: %w | stderr: %s", err, stderr.String())
	}
	
	return allStreams, nil
}
