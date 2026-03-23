package scraper

// Reddit scraper — stubbed for Day 1, enable when credentials are ready
//
// Setup:
//   1. Go to https://www.reddit.com/prefs/apps and create a "script" app
//   2. Set REDDIT_CLIENT_ID, REDDIT_CLIENT_SECRET, REDDIT_USER_AGENT in .env
//
// This file is intentionally left as a stub to keep the build compiling.
// Uncomment the full implementation and add to manager.go scrapers slice when ready.

/*
import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/kaushikduvur/baltimore-sports-scraper/backend/internal/models"
)

const redditBaseURL = "https://oauth.reddit.com"

type RedditScraper struct {
	client      *http.Client
	accessToken string
	tokenExpiry time.Time
	clientID    string
	clientSecret string
	userAgent   string
}

func NewRedditScraper() *RedditScraper {
	return &RedditScraper{
		client:       &http.Client{Timeout: 15 * time.Second},
		clientID:     os.Getenv("REDDIT_CLIENT_ID"),
		clientSecret: os.Getenv("REDDIT_CLIENT_SECRET"),
		userAgent:    os.Getenv("REDDIT_USER_AGENT"),
	}
}

func (r *RedditScraper) Name() string            { return "Reddit" }
func (r *RedditScraper) Interval() time.Duration { return 5 * time.Minute }

func (r *RedditScraper) Fetch(ctx context.Context) ([]models.Article, error) {
	if err := r.ensureToken(ctx); err != nil {
		return nil, err
	}

	subreddits := []struct {
		sub  string
		team models.Team
	}{
		{"orioles", models.TeamOrioles},
		{"ravens", models.TeamRavens},
		{"baseball", models.TeamOrioles},
		{"nfl", models.TeamRavens},
	}

	var articles []models.Article
	for _, sr := range subreddits {
		items, err := r.fetchSubreddit(ctx, sr.sub, sr.team)
		if err != nil {
			continue
		}
		articles = append(articles, items...)
	}
	return articles, nil
}
*/
