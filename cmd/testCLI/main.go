package main

import (
	"fmt"
	"log"
	"os"

	"github.com/coopons/livestream_scraper/internal/api"
)

func main() {
	clientID := os.Getenv("TWITCH_CLIENT_ID")
	clientSecret := os.Getenv("TWITCH_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		log.Fatal("Please set TWITCH_CLIENT_ID and TWITCH_CLIENT_SECRET")
	}

	token, err := api.GetAppAccessToken(clientID, clientSecret)
	if err != nil {
		log.Fatalf("Failed to get access token: %v", err)
	}

	// streams, err := api.GetLiveStreams(clientID, token, 10)
	// if err != nil {
	// 	log.Fatal("Failed to get live streams:", err)
	// }
	streams, err := api.GetAllLiveStreams(clientID, token, 250)
	if err != nil {
		log.Fatal("Error getting all streams:", err)
	}

	fmt.Println("Total streams collected:", len(streams))

	var streamTitles []string
	for _, s := range streams {
		streamTitles = append(streamTitles, fmt.Sprintf("%s - %s, (%d viewers)", s.UserName, s.Title, s.ViewerCount))
	}

	fmt.Println("Live Streams:")
	for _, title := range streamTitles {
		fmt.Println(title)
	}
	if len(streams)==0 {
		fmt.Println("No streams available!")
	}
}