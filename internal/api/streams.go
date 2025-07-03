package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Stream struct {
	ID 		string `json:"id"`
	UserID 	string `json:"user_id"`
	UserName string `json:"user_name"`
	Title 	string `json:"title"`
	GameID 	string `json:"game_id"`
	Language string `json:"language"`
	ViewerCount int    `json:"viewer_count"`
	StartedAt string `json:"started_at"`
	ThumbnailURL string `json:"thumbnail_url"`
}

type GetStreamsResponse struct {
	Data []Stream `json:"data"`
}

func GetLiveStreams(clientID, accessToken string, first int) ([]Stream, error) {
	url := fmt.Sprintf("https://api.twitch.tv/helix/streams?first=%d", first)
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