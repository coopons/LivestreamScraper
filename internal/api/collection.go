package api

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/coopons/livestream_scraper/internal/db"
	"github.com/coopons/livestream_scraper/internal/scraper"
)

var (
	ticker			 *time.Ticker
	tickerStop	 	 chan bool
	nextRunTime	 	 time.Time
	tickerInterval	 time.Duration
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
	// Prevent collection running too often during testing
	latestSnapshot, err := db.GetLatestSnapshotTime()
	if err != nil {
		log.Println("Error getting latest snapshot:", err)
	} else {
		if time.Since(latestSnapshot) < 4*time.Minute {
			log.Println("Skipping collection, last snapshot was less than 4 minutes ago")
			return nil
		}
	}

	// --- TWITCH ---
	token, err := GetAppAccessToken(clientID, clientSecret)
	if err != nil {
		return err
	}

	twitchStreams, err := GetAllLiveStreams(clientID, token, 200) // Only get top 200 streams to start out
	if err != nil {
		return err
	}

	for _, s := range twitchStreams {
		if err := db.SaveStream(s, "twitch"); err != nil {
			log.Println("SaveStream (twitch) error:", err)
		}
		if err := db.SaveSnapshot(s); err != nil {
			log.Println("SaveSnapshot (twitch) error:", err)
		}
	}

	// --- YOUTUBE ---
	ytStreams, err := scraper.ScrapeYoutubeLivestreams()
	if err != nil {
		log.Println("ScrapeYoutubeLivestreams error:", err)
	} else {
		for _, s := range ytStreams {
			if err := db.SaveStream(s, "youtube"); err != nil {
				log.Println("SaveStream (youtube) error:", err)
			}
			if err := db.SaveSnapshot(s); err != nil {
				log.Println("SaveSnapshot (youtube) error:", err)
			}
		}
	}

	log.Printf("Collected %d Twitch streams and %d YouTube streams\n", len(twitchStreams), len(ytStreams))

	return nil
}

func StopCollector() {
	if tickerStop != nil {
		select {
		case tickerStop <- true:
		default:
		}
	}
	collectorRunning = false
}

func ControlHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Query().Get("action") {
	case "stop":
		StopCollector()
	case "start":
		clientID := os.Getenv("TWITCH_CLIENT_ID")
		clientSecret := os.Getenv("TWITCH_CLIENT_SECRET")
		StartCollector(clientID, clientSecret, 5*time.Minute)
	default: 
		log.Println("Unknown action:", r.URL.Query().Get("action"))
	}
	w.WriteHeader(http.StatusOK)
}
