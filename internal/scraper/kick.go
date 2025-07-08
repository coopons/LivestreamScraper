package scraper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/coopons/livestream_scraper/internal/model"
)

type KickResponse struct {
	Data []model.KickStream `json:"data"`
	Message string            `json:"message"`
}

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
	baseURL := "https://api.kick.com/public/v1/livestreams"
	u, err := url.Parse(baseURL)
	if err != nil {
		panic(err)
	}

	q := u.Query()
	q.Set("limit", "100")
	q.Set("sort", "viewer_count")
	u.RawQuery = q.Encode()

	
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")
	
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch kick streams: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("kick API error (%d): %s", resp.StatusCode, resp.Status)
	}

	var respData KickResponse
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, fmt.Errorf("error decoding kick streams JSON: %w", err)
	}
	
	allStreams := make([]model.Stream, 0, len(respData.Data))
	for _, ks := range respData.Data {
		allStreams = append(allStreams, ks.ToModelStream())
	}

	return allStreams, nil
}
