package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	goredis "github.com/redis/go-redis/v9"

	"github.com/kaushikduvur/baltimore-sports-scraper/backend/internal/models"
)

const (
	StreamKey      = "articles:stream"
	StreamMaxLen   = 1000 // Keep last 1000 events
	ConsumerGroup  = "sse-consumers"
)

// NewClient creates a Redis client from a URL
func NewClient(redisURL string) *goredis.Client {
	opts, err := goredis.ParseURL(redisURL)
	if err != nil {
		// Fallback to default local Redis
		log.Printf("Failed to parse REDIS_URL, using default localhost:6379: %v", err)
		opts = &goredis.Options{Addr: "localhost:6379"}
	}
	return goredis.NewClient(opts)
}

// PublishArticle pushes a new article event to the Redis Stream
func PublishArticle(ctx context.Context, rdb *goredis.Client, article models.Article) error {
	payload, err := json.Marshal(article)
	if err != nil {
		return fmt.Errorf("failed to marshal article: %w", err)
	}

	return rdb.XAdd(ctx, &goredis.XAddArgs{
		Stream: StreamKey,
		MaxLen: StreamMaxLen,
		Approx: true,
		Values: map[string]interface{}{
			"article": string(payload),
			"team":    string(article.Team),
			"source":  string(article.Source),
		},
	}).Err()
}

// Subscribe reads new messages from the stream starting from the given offset.
// Pass "$" to only receive messages published after Subscribe is called.
// Pass "0" to replay from the beginning.
func Subscribe(ctx context.Context, rdb *goredis.Client, lastID string, ch chan<- models.Article) {
	if lastID == "" {
		lastID = "$"
	}
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		streams, err := rdb.XRead(ctx, &goredis.XReadArgs{
			Streams: []string{StreamKey, lastID},
			Count:   10,
			Block:   5 * time.Second,
		}).Result()

		if err == goredis.Nil {
			continue // timeout with no messages — loop
		}
		if err != nil {
			log.Printf("Redis stream read error: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		for _, stream := range streams {
			for _, msg := range stream.Messages {
				lastID = msg.ID
				raw, ok := msg.Values["article"].(string)
				if !ok {
					continue
				}
				var article models.Article
				if err := json.Unmarshal([]byte(raw), &article); err != nil {
					log.Printf("Failed to unmarshal article from stream: %v", err)
					continue
				}
				select {
				case ch <- article:
				case <-ctx.Done():
					return
				}
			}
		}
	}
}
