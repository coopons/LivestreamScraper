package web

import (
	"log"
	"time"

	"github.com/coopons/livestream_scraper/internal/api"
	"github.com/coopons/livestream_scraper/internal/db"
)

var (
	ticker		 *time.Ticker
	tickerStop	 chan bool
	tickerPaused bool
	nextRunTime	 time.Time
	tickerInterval time.Duration
)

func StartCollector(clientID, clientSecret string, interval time.Duration) {
	if ticker != nil {
		ticker.Stop()
	}
	tickerStop = make(chan bool)
	ticker = time.NewTicker(interval)
	tickerInterval = interval
	tickerPaused = false

	go func() {
		for {
			select {
			case <-ticker.C:
				if tickerPaused {
					continue
				}
				nextRunTime = time.Now().Add(tickerInterval)
				err := runCollection(clientID, clientSecret)
				if err != nil {
					log.Println("Start collection error:", err)
				}
			case <-tickerStop:
				ticker.Stop()
				log.Println("Collector stopped.")
				return
			}
		}
	}()
}

func runCollection(clientID, clientSecret string) error {
	log.Println("Running Collection...")
	token, err := api.GetAppAccessToken(clientID, clientSecret)
	if err != nil {
		return err
	}

	streams, err := api.GetAllLiveStreams(clientID, token, 200) // Only get top 200 streams to start out
	if err != nil {
		return err
	}

	for _, s := range streams {
		if err := db.SaveStream(s, "twitch"); err != nil {
			log.Println("SaveStream error:", err)
		}
		if err := db.SaveSnapshot(s); err != nil {
			log.Println("SaveSnapshot error:", err)
		}
	}
	return nil
}

func pauseCollector() {
	tickerPaused = true
	log.Println("Collector paused.")
}

func resumeCollector() {
	tickerPaused = false
	log.Println("Collector resumed.")
	nextRunTime = time.Now().Add(tickerInterval)
}

func stopCollector() {
	if tickerStop != nil {
		tickerStop <- true
	}
}