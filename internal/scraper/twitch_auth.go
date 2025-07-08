package scraper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"golang.org/x/oauth2"
)

var ( 
	twitchToken 			*oauth2.Token
	twitchTokenMutex		sync.Mutex
)

func GetTwitchCachedToken(clientID, clientSecret string) (string, error) {
	twitchTokenMutex.Lock()
	defer twitchTokenMutex.Unlock()

	if twitchToken != nil && twitchToken.Valid() {
		return twitchToken.AccessToken, nil
	}
	return fetchTwitchNewToken(clientID, clientSecret)
}

func fetchTwitchNewToken(clientID, clientSecret string) (string, error) {
	resp, err := http.PostForm("https://id.twitch.tv/oauth2/token", url.Values{
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"grant_type":    {"client_credentials"},
	})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	type tokenResponse struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("twitch token error (%d): %s", resp.StatusCode, string(body))
	}

	var tr tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return "", err
	}

	twitchToken = &oauth2.Token{
		AccessToken: tr.AccessToken,
		Expiry:      time.Now().Add(time.Duration(tr.ExpiresIn) * time.Second),
	}

	return twitchToken.AccessToken, nil
}