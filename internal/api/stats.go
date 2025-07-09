package api

import (
	"context"

	"github.com/coopons/livestream_scraper/internal/db"
)

type StatsPageData struct {
	AverageDuration    []AverageDurationEntry
	Peak30             int
	Peak30Streamer	   string
	PeakAllTime        int
	PeakAllStreamer	   string
	PopularTimes       []PopularTimeEntry
	TopCategories      []TopCategoryEntry
	PeakHourComparison []PeakHourComparisonEntry
}

type AverageDurationEntry struct {
	Category    string
	Day         string
	AvgDuration int
}

type PopularTimeEntry struct {
	Platform   string
	Hour       int
	AvgViewers int
}

type TopCategoryEntry struct {
	Platform   string
	Category   string
	AvgViewers int
}

type PeakHourComparisonEntry struct {
	Hour       int
	Platform   string
	AvgViewers int
}

func GetStatsPageData() (StatsPageData, error) {
	var data StatsPageData
	// Average stream duration per category and day
	avgDurationQuery := `
		SELECT game, TO_CHAR(timestamp, 'Day') AS day, ROUND(AVG(duration_minutes)) AS avg_duration
		FROM (
			SELECT s.stream_id, game, DATE_TRUNC('minute', MAX(timestamp) - MIN(timestamp)) AS duration,
			EXTRACT(EPOCH FROM MAX(timestamp) - MIN(timestamp)) / 60 AS duration_minutes,
			MIN(timestamp) AS timestamp
			FROM stream_snapshots s
			JOIN streams ON streams.stream_id = s.stream_id
			GROUP BY s.stream_id, game
		) durations
		GROUP BY game, day
		ORDER BY day, avg_duration DESC`
	rows, err := db.Pool.Query(context.Background(), avgDurationQuery)
	if err != nil {
		return data, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry AverageDurationEntry
		err := rows.Scan(&entry.Category, &entry.Day, &entry.AvgDuration)
		if err != nil {
			return data, err
		}
		data.AverageDuration = append(data.AverageDuration, entry)
	}

	// Peak viewer numbers (last 30 days and all time) with streamer name
	err = db.Pool.QueryRow(context.Background(), `
		SELECT viewer_count, streamer_name
		FROM stream_snapshots
		JOIN streams ON stream_snapshots.stream_id = streams.stream_id
		WHERE timestamp >= NOW() - INTERVAL '30 days'
		ORDER BY viewer_count DESC
		LIMIT 1`).Scan(&data.Peak30, &data.Peak30Streamer)
	if err != nil {
		return data, err
	}

	err = db.Pool.QueryRow(context.Background(), `
		SELECT viewer_count, streamer_name
		FROM stream_snapshots
		JOIN streams ON stream_snapshots.stream_id = streams.stream_id
		ORDER BY viewer_count DESC
		LIMIT 1`).Scan(&data.PeakAllTime, &data.PeakAllStreamer)
	if err != nil {
		return data, err
	}


	// Most popular times by platform (hour with highest average viewers)
	popularTimesQuery := `
		SELECT platform, EXTRACT(HOUR FROM timestamp)::int AS hour, ROUND(AVG(viewer_count)) AS avg_viewers
		FROM stream_snapshots
		JOIN streams ON stream_snapshots.stream_id = streams.stream_id
		GROUP BY platform, hour
		ORDER BY platform, avg_viewers DESC`
	rows, err = db.Pool.Query(context.Background(), popularTimesQuery)
	if err != nil {
		return data, err
	}
	defer rows.Close()

	// Use a map to keep highest hour per platform
	bestHour := make(map[string]PopularTimeEntry)
	for rows.Next() {
		var e PopularTimeEntry
		err := rows.Scan(&e.Platform, &e.Hour, &e.AvgViewers)
		if err != nil {
			return data, err
		}
		if _, exists := bestHour[e.Platform]; !exists {
			bestHour[e.Platform] = e
		}
	}
	for _, v := range bestHour {
		data.PopularTimes = append(data.PopularTimes, v)
	}

	// Top categories per platform (by average viewer count)
	topCategoryQuery := `
		SELECT platform, game, ROUND(AVG(viewer_count)) AS avg_viewers
		FROM stream_snapshots
		JOIN streams ON stream_snapshots.stream_id = streams.stream_id
		GROUP BY platform, game
		ORDER BY platform, avg_viewers DESC`
	rows, err = db.Pool.Query(context.Background(), topCategoryQuery)
	if err != nil {
		return data, err
	}
	defer rows.Close()

	topCategorySeen := make(map[string]TopCategoryEntry)
	for rows.Next() {
		var e TopCategoryEntry
		err := rows.Scan(&e.Platform, &e.Category, &e.AvgViewers)
		if err != nil {
			return data, err
		}
		if _, ok := topCategorySeen[e.Platform]; !ok {
			topCategorySeen[e.Platform] = e
		}
	}
	for _, v := range topCategorySeen {
		data.TopCategories = append(data.TopCategories, v)
	}

	// Viewer distribution at peak hours
	peakHourComparisonQuery := `
		SELECT hour, platform, avg_viewers
			FROM (
				SELECT 
					EXTRACT(HOUR FROM stream_snapshots.timestamp)::int AS hour,
					platform,
					ROUND(AVG(viewer_count)) AS avg_viewers,
					ROW_NUMBER() OVER (
						PARTITION BY EXTRACT(HOUR FROM stream_snapshots.timestamp), platform
						ORDER BY AVG(viewer_count) DESC
					) AS rank
				FROM stream_snapshots
				JOIN streams ON stream_snapshots.stream_id = streams.stream_id
				GROUP BY EXTRACT(HOUR FROM stream_snapshots.timestamp), platform
			) ranked
			WHERE rank = 1;`
	rows, err = db.Pool.Query(context.Background(), peakHourComparisonQuery)
	if err != nil {
		return data, err
	}
	defer rows.Close()

	for rows.Next() {
		var e PeakHourComparisonEntry
		err := rows.Scan(&e.Hour, &e.Platform, &e.AvgViewers)
		if err != nil {
			return data, err
		}
		data.PeakHourComparison = append(data.PeakHourComparison, e)
	}

	return data, nil
}
