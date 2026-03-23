package scraper

// Bluesky AT Protocol scraper — stubbed for Day 1
//
// Setup:
//   1. Create a Bluesky account at https://bsky.app
//   2. Generate an app password: Settings → App Passwords
//   3. Set BLUESKY_IDENTIFIER (your handle) and BLUESKY_APP_PASSWORD in .env
//
// The AT Protocol (atproto.com) is fully open and free — no rate limits for read-only access.
// Uses the app.bsky.feed.searchPosts lexicon to find posts about Orioles/Ravens.

/*
import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/kaushikduvur/baltimore-sports-scraper/backend/internal/models"
)

const bskyBaseURL = "https://bsky.social/xrpc"

type BlueskyScraper struct {
	client     *http.Client
	accessJWT  string
	did        string
	identifier string
	appPassword string
}

func NewBlueskyScraper() *BlueskyScraper {
	return &BlueskyScraper{
		client:      &http.Client{Timeout: 15 * time.Second},
		identifier:  os.Getenv("BLUESKY_IDENTIFIER"),
		appPassword: os.Getenv("BLUESKY_APP_PASSWORD"),
	}
}

func (b *BlueskyScraper) Name() string            { return "Bluesky" }
func (b *BlueskyScraper) Interval() time.Duration { return 3 * time.Minute }

func (b *BlueskyScraper) Fetch(ctx context.Context) ([]models.Article, error) {
	if err := b.authenticate(ctx); err != nil {
		return nil, err
	}
	// Search posts about Baltimore teams
	queries := []struct {
		q    string
		team models.Team
	}{
		{"Baltimore Orioles", models.TeamOrioles},
		{"Baltimore Ravens", models.TeamRavens},
		{"#Orioles", models.TeamOrioles},
		{"#Ravens", models.TeamRavens},
	}
	var articles []models.Article
	for _, q := range queries {
		items, _ := b.searchPosts(ctx, q.q, q.team)
		articles = append(articles, items...)
	}
	return articles, nil
}
*/
