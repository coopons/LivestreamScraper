package scraper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/coopons/livestream_scraper/internal/model"
)

type KickScraper struct {
	ClientID	 string
	ClientSecret string
	Token		 string
}

func (k *KickScraper) GetLiveStreams(limit int) ([]model.Stream, error) {
	token, err := GetKickCachedToken(k.ClientID, k.ClientSecret)
	if err != nil {
		return nil, err
	}
	return getKickStreams(token, limit) // Set to top 200 streams for now
}

func (k *KickScraper) Platform() string {
	return "kick"
}

// Fetches the specified number of kick streams sorted by viewer count 
func getKickStreams(accessToken string, max int) ([]model.Stream, error) {
	url := fmt.Sprintf("https://api.kick.com/api/v1/livestreams?sort=viewer_count&limit=%d", max)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")
	
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch kick streams: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("kick API error (%d): %s", resp.StatusCode, resp.Status)
	}

	var kickStreams []model.KickStream
	if err := json.NewDecoder(resp.Body).Decode(&kickStreams); err != nil {
		return nil, fmt.Errorf("error decoding kick streams JSON: %w", err)
	}
	
	allStreams := make([]model.Stream, 0, len(kickStreams))
	for _, ks := range kickStreams {
		allStreams = append(allStreams, ks.ToModelStream())
	}

	return allStreams, nil
}
