package scraper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

var ( 
	twitchOAuthConfig  	*oauth2.Config
	token 				*oauth2.Token
)

func InitTwitchOAuth(clientID, clientSecret, redirectURL string) {
	twitchOAuthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://id.twitch.tv/oauth2/authorize",
			TokenURL: "https://id.twitch.tv/oauth2/token",
		},
		RedirectURL: "http://localhost:",
		Scopes: 	 []string{"user:read:email", "user:read:follows"},
	}
}

func GetAppAccessToken(clientID, clientSecret string) (string, error) {
	url := fmt.Sprintf("https://id.twitch.tv/oauth2/token?client_id=%s&client_secret=%s&grant_type=client_credentials", clientID, clientSecret)
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	type tokenResponse struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}

	var tr tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return "", err
	}

	token = &oauth2.Token{
		AccessToken: tr.AccessToken,
		Expiry:    time.Now().Add(time.Duration(tr.ExpiresIn) * time.Second),
	}
	return tr.AccessToken, nil
}