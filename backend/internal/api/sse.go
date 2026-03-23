package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	redisClient "github.com/kaushikduvur/baltimore-sports-scraper/backend/internal/redis"
	"github.com/kaushikduvur/baltimore-sports-scraper/backend/internal/models"
)

// SSEStream godoc
// GET /api/stream?team=orioles
//
// Opens a Server-Sent Events connection. Clients receive new articles in real time
// as they are scraped. The connection stays open indefinitely.
//
// Event format:
//   event: article
//   data: {"id":"...","title":"...","team":"orioles",...}
func (h *Handler) SSEStream(w http.ResponseWriter, r *http.Request) {
	// Verify the client supports SSE (flushing)
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	teamFilter := r.URL.Query().Get("team") // optional filter

	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no") // Disable Nginx buffering

	// Send a connected event immediately
	fmt.Fprintf(w, "event: connected\ndata: {\"status\":\"connected\"}\n\n")
	flusher.Flush()

	// Subscribe to Redis stream
	ch := make(chan models.Article, 10)
	ctx := r.Context()

	go redisClient.Subscribe(ctx, h.rdb, "$", ch)

	// Keep-alive ticker — browsers disconnect if no data arrives for ~30s
	keepAlive := time.NewTicker(20 * time.Second)
	defer keepAlive.Stop()

	log.Printf("SSE client connected: %s (team filter: %q)", r.RemoteAddr, teamFilter)

	for {
		select {
		case <-ctx.Done():
			log.Printf("SSE client disconnected: %s", r.RemoteAddr)
			return

		case article := <-ch:
			// Apply optional team filter
			if teamFilter != "" && string(article.Team) != teamFilter && article.Team != models.TeamBoth {
				continue
			}

			data, err := json.Marshal(article)
			if err != nil {
				continue
			}

			fmt.Fprintf(w, "event: article\ndata: %s\n\n", data)
			flusher.Flush()

		case <-keepAlive.C:
			// SSE comment keeps connection alive without triggering client events
			fmt.Fprintf(w, ": keepalive\n\n")
			flusher.Flush()
		}
	}
}
