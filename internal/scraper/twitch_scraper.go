package scraper

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

func ScrapeTwitchBrowse() error {
	c := colly.NewCollector(
		colly.AllowedDomains("www.twitch.tv", "twitch.tv"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"),
	)
	
	c.OnHTML("div[data-test-selector='DirectoryCard']", func(e *colly.HTMLElement) {
		streamer := e.ChildText("p.ScThumbnailCardInfo-username")
		title := e.ChildText("h3.ScThumbnailCardInfo-title")
		viewer_count := e.ChildText("p.ScThumbnailCardInfo-viewers")
		category := e.ChildText("p.ScThumbnailCardInfo-gameName")

		fmt.Printf("Streamer: %s, Title: %s, Category: %s, Viewers: %s\n", streamer, title, category, viewer_count)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	return c.Visit("https://www.twitch.tv/directory")
}