package api

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/coopons/livestream_scraper/internal/db"
	"github.com/coopons/livestream_scraper/internal/model"
	"github.com/coopons/livestream_scraper/internal/scraper"
)

var (
	ticker			 *time.Ticker
	tickerStop	 	 chan bool
	nextRunTime	 	 time.Time
	tickerInterval	 time.Duration
	collectorRunning bool
)

func StartCollector(twitchClientID, twitchClientSecret, kickClientID, kickClientSecret string, interval time.Duration) {
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

	collectors := []scraper.StreamCollector{
		&scraper.TwitchScraper{ClientID: twitchClientID, ClientSecret: twitchClientSecret},
		&scraper.YoutubeScraper{},
		&scraper.KickScraper{ClientID: kickClientID, ClientSecret: kickClientSecret},
	}

	// Attempt collection at program start
	go func() {
		if err := runCollection(collectors); err != nil {
			log.Println("Initial collection error:", err)
		} 
	}()

	// Start the ticker to run collection periodically
	go func() {
		for {
			select {
			case <-ticker.C:
				nextRunTime = time.Now().Add(tickerInterval)
				err := runCollection(collectors)
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

func runCollection(collectors []scraper.StreamCollector) error {
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

	type result struct {
		streams  []model.Stream
		platform string
		duration time.Duration
		err		 error
	}

	resultsCh := make(chan result)

	for _, collector := range collectors {
		go func(c scraper.StreamCollector) {
			start := time.Now()
			streams, err := c.GetLiveStreams(200)
			duration := time.Since(start)
			
			resultsCh <- result{
				streams: streams,
				platform: c.Platform(),
				duration: duration,
				err: err,
			}
		}(collector)
	}

	for i := 0; i < len(collectors); i++ {
		res := <-resultsCh
		if res.err != nil {
			log.Printf("Error collecting %s streams: %v\n", res.platform, res.err)
			continue
		}
		for _, s := range res.streams {
			if err := db.SaveStream(s, res.platform); err != nil {
				log.Printf("SaveStream (%s) error: %v\n", res.platform, err)
			}
			if err := db.SaveSnapshot(s); err != nil {
				log.Printf("SaveSnapshot (%s) error: %v\n", res.platform, err)
			}
		}
		log.Printf("Collected %d %s streams in %s\n", len(res.streams), res.platform, res.duration)
	}

	lastRunTime = time.Now()

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
		twitchClientID := os.Getenv("TWITCH_CLIENT_ID")
		twitchClientSecret := os.Getenv("TWITCH_CLIENT_SECRET")
		kickClientID := os.Getenv("KICK_CLIENT_ID")
		kickClientSecret := os.Getenv("KICK_CLIENT_SECRET")
		StartCollector(twitchClientID, twitchClientSecret, kickClientID, kickClientSecret, 10*time.Minute)
	default: 
		log.Println("Unknown action:", r.URL.Query().Get("action"))
	}
	w.WriteHeader(http.StatusOK)
}
