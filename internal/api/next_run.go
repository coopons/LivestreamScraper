package api

import (
	"encoding/json"
	"net/http"
	"time"
)

func NextRunHandler(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		NextRun string `json:"next_run"`
		Running bool   `json:"running"`
	}

	if nextRunTime.IsZero() {
		nextRunTime = time.Now().Add(tickerInterval)
	}

	resp := Response{
		NextRun: nextRunTime.Format(time.RFC3339),
		Running: collectorRunning,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}