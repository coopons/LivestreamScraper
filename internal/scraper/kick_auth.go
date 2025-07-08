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
	kickToken 			*oauth2.Token
	kickTokenMutex		sync.Mutex
)

func GetKickCachedToken(clientID, clientSecret string) (string, error) {
	kickTokenMutex.Lock()
	defer kickTokenMutex.Unlock()

	if kickToken != nil && kickToken.Valid() {
		return kickToken.AccessToken, nil
	}
	return fetchKickNewToken(clientID, clientSecret)
}

func fetchKickNewToken(clientID, clientSecret string) (string, error) {
	resp, err := http.PostForm("https://id.kick.com/oauth/token", url.Values{
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
		return "", fmt.Errorf("kick token error (%d): %s", resp.StatusCode, string(body))
	}

	var tr tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return "", err
	}

	kickToken = &oauth2.Token{
		AccessToken: tr.AccessToken,
		Expiry:      time.Now().Add(time.Duration(tr.ExpiresIn) * time.Second),
	}

	return kickToken.AccessToken, nil
}