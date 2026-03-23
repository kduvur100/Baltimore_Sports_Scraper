package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/kaushikduvur/baltimore-sports-scraper/backend/internal/models"
)

// GDELT DOC 2.0 API — free, no API key required
// Docs: https://blog.gdeltproject.org/gdelt-doc-2-0-api-debuts/
const gdeltBaseURL = "https://api.gdeltproject.org/api/v2/doc/doc"

type gdeltResponse struct {
	Articles []gdeltArticle `json:"articles"`
}

type gdeltArticle struct {
	URL        string  `json:"url"`
	Title      string  `json:"title"`
	Seendate   string  `json:"seendate"`
	Domain     string  `json:"domain"`
	Language   string  `json:"language"`
	ToneAvg    float64 `json:"tone"`
	SocialImage string `json:"socialimage"`
}

type GDELTScraper struct {
	client *http.Client
}

func NewGDELTScraper() *GDELTScraper {
	return &GDELTScraper{
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (g *GDELTScraper) Name() string            { return "GDELT" }
func (g *GDELTScraper) Interval() time.Duration { return 15 * time.Minute }

func (g *GDELTScraper) Fetch(ctx context.Context) ([]models.Article, error) {
	var articles []models.Article

	// Fetch for each team separately
	queries := []struct {
		query string
		team  models.Team
	}{
		{`"Baltimore Orioles" OR "Orioles baseball"`, models.TeamOrioles},
		{`"Baltimore Ravens" OR "Ravens NFL"`, models.TeamRavens},
	}

	for _, q := range queries {
		items, err := g.fetchQuery(ctx, q.query, q.team)
		if err != nil {
			// Log but don't fail entirely
			continue
		}
		articles = append(articles, items...)
	}

	return articles, nil
}

func (g *GDELTScraper) fetchQuery(ctx context.Context, query string, team models.Team) ([]models.Article, error) {
	params := url.Values{}
	params.Set("query", query)
	params.Set("mode", "artlist")
	params.Set("maxrecords", "25")
	params.Set("format", "json")
	params.Set("timespan", "1d") // Last 24 hours only
	params.Set("sort", "datedesc")

	reqURL := fmt.Sprintf("%s?%s", gdeltBaseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "BaltimoreSportsScraper/1.0")

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GDELT API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result gdeltResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse GDELT response: %w", err)
	}

	var articles []models.Article
	for _, a := range result.Articles {
		if a.Language != "English" {
			continue
		}

		pub, err := time.Parse("20060102T150405Z", a.Seendate)
		if err != nil {
			pub = time.Now()
		}

		// Normalize tone (-100 to +100) to 0–1 sentiment scale
		sentiment := (a.ToneAvg + 100) / 200

		articles = append(articles, models.Article{
			Title:          a.Title,
			URL:            a.URL,
			Summary:        fmt.Sprintf("via %s", a.Domain),
			Source:         models.SourceGDELT,
			Team:           team,
			ImageURL:       a.SocialImage,
			SentimentScore: sentiment,
			PublishedAt:    pub,
		})
	}

	return articles, nil
}
