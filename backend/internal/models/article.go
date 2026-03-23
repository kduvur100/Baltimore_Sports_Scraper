package models

import "time"

// Source represents which data pipeline an article came from
type Source string

const (
	SourceGDELT    Source = "gdelt"
	SourceReddit   Source = "reddit"
	SourceBluesky  Source = "bluesky"
	SourceRSS      Source = "rss"
	SourceYouTube  Source = "youtube"
)

// Team tag
type Team string

const (
	TeamOrioles Team = "orioles"
	TeamRavens  Team = "ravens"
	TeamBoth    Team = "both"
)

// Article is the unified content model used across all sources
type Article struct {
	ID          string    `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	URL         string    `json:"url" db:"url"`
	Summary     string    `json:"summary" db:"summary"`
	Source      Source    `json:"source" db:"source"`
	Team        Team      `json:"team" db:"team"`
	Author      string    `json:"author,omitempty" db:"author"`
	ImageURL    string    `json:"image_url,omitempty" db:"image_url"`
	SentimentScore float64 `json:"sentiment_score,omitempty" db:"sentiment_score"`
	PublishedAt time.Time `json:"published_at" db:"published_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// ArticleEvent is published to Redis Streams when a new article arrives
type ArticleEvent struct {
	Article   Article `json:"article"`
	EventType string  `json:"event_type"` // "new_article"
}
