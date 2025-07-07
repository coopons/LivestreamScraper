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
	token 			*oauth2.Token
	tokenMutex		sync.Mutex
)

func GetCachedToken(clientID, clientSecret string) (string, error) {
	tokenMutex.Lock()
	defer tokenMutex.Unlock()

	if token != nil && token.Valid() {
		return token.AccessToken, nil
	}
	return fetchNewToken(clientID, clientSecret)
}

func fetchNewToken(clientID, clientSecret string) (string, error) {
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

	token = &oauth2.Token{
		AccessToken: tr.AccessToken,
		Expiry:      time.Now().Add(time.Duration(tr.ExpiresIn) * time.Second),
	}

	return tr.AccessToken, nil
}