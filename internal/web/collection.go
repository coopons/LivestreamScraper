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
	nextRunTime	 time.Time
	tickerInterval time.Duration
	collectorRunning bool
)

func StartCollector(clientID, clientSecret string, interval time.Duration) {
	if collectorRunning {
		log.Println("Collector is already running.")
		return
	}
	if ticker != nil {
		ticker.Stop()
	}
	if tickerStop != nil {
		select {
		case tickerStop <- true:
		default:
		}
	}

	tickerStop = make(chan bool)
	ticker = time.NewTicker(interval)
	tickerInterval = interval
	nextRunTime = time.Now().Add(tickerInterval)
	collectorRunning = true

	go func() {
		for {
			select {
			case <-ticker.C:
				nextRunTime = time.Now().Add(tickerInterval)
				err := runCollection(clientID, clientSecret)
				if err != nil {
					log.Println("Start collection error:", err)
				}
			case <-tickerStop:
				ticker.Stop()
				collectorRunning = false
 				return
			}
		}
	}()
}

func runCollection(clientID, clientSecret string) error {
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

func stopCollector() {
	if tickerStop != nil {
		select {
		case tickerStop <- true:
		default:
		}
	}
	collectorRunning = false
}