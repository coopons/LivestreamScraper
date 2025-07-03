package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Data parsed from the twitch api curl
type Stream struct {
	ID 				string `json:"id"`
	UserID 			string `json:"user_id"`
	UserName 		string `json:"user_name"`
	Title 			string `json:"title"`
	GameID 			string `json:"game_id"`
	GameName 		string `json:"game_name"`
	Language 		string `json:"language"`
	ViewerCount 	int    `json:"viewer_count"`
	StartedAt 		string `json:"started_at"`
	ThumbnailURL 	string `json:"thumbnail_url"`
	IsMature 		bool   `json:"is_mature"`
}

type GetStreamsResponse struct {
	Data []Stream `json:"data"`
	Pagination struct {
		Cursor string `json:"cursor"`
	} `json:"pagination"`
}

// Gets specified number of streams from twitch API
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

// Continually requests streams from Twitch API until all streams are retrieved
func GetAllLiveStreams(clientID, accessToken string, max int) ([]Stream, error) {
	var allStreams []Stream
	cursor := ""
	pageSize := 100
	client := &http.Client{}

	for {
		if len(allStreams) >= max {
			break
		}
		url := fmt.Sprintf("https://api.twitch.tv/helix/streams?first=%d", pageSize)
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

