package scraper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/coopons/livestream_scraper/internal/model"
)

type GetStreamsResponse struct {
	Data []model.Stream `json:"data"`
	Pagination struct {
		Cursor string `json:"cursor"`
	} `json:"pagination"`
}

type TwitchScraper struct {
	ClientID	 string
	ClientSecret string
	Token		 string
}

func (t *TwitchScraper) GetLiveStreams(limit int) ([]model.Stream, error) {
	if t.ClientSecret == "" {
		token, err := GetAppAccessToken(t.ClientID, t.ClientSecret)
		if err != nil {
			return nil, err
		}
		t.Token = token
	}
	return GetAllLiveStreams(t.ClientID, t.Token, limit) // Set to top 200 streams for now
}

func (t *TwitchScraper) Platform() string {
	return "twitch"
}

// Recursively fetches the top twitch streams until max is reached
func GetAllLiveStreams(clientID, accessToken string, max int) ([]model.Stream, error) {
	var allStreams []model.Stream
	cursor := ""
	pageSize := 100
	client := &http.Client{}
	
	for {
		if len(allStreams) >= max {
			break
		}
		url := fmt.Sprintf("https://api.twitch.tv/helix/streams?type=live&first=%d", pageSize)
		if cursor != "" {
			url += "&after=" + cursor
		}
		
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		
		req.Header.Set("Client-ID", clientID)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("failed to get streams: %s", resp.Status)
		}
		
		var response GetStreamsResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, err
		}
		
		for _, stream := range response.Data {
			stream.ThumbnailURL = strings.ReplaceAll(stream.ThumbnailURL, "{width}", "320")
			stream.ThumbnailURL = strings.ReplaceAll(stream.ThumbnailURL, "{height}", "180")
			allStreams = append(allStreams, stream)
			
			if len(allStreams) >= max {
				break
			}
		}
		
		if response.Pagination.Cursor == "" {
			break
		}
		
		cursor = response.Pagination.Cursor
	}
	return allStreams, nil
}

// Gets specified number of streams from twitch API
// --- NO LONGER USED ---
func GetLiveStreams(clientID, accessToken string, first int) ([]model.Stream, error) {
	url := fmt.Sprintf("https://api.twitch.tv/helix/streams?type=live&first=%d", first)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Client-ID", clientID)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get streams: %s", resp.Status)
	}

	var response GetStreamsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	for i := range response.Data {
		response.Data[i].ThumbnailURL = strings.ReplaceAll(response.Data[i].ThumbnailURL, "{width}", "320")
		response.Data[i].ThumbnailURL = strings.ReplaceAll(response.Data[i].ThumbnailURL, "{height}", "180")
	}

	return response.Data, nil
}
