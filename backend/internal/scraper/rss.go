package scraper

import (
	"context"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"

	"github.com/kaushikduvur/baltimore-sports-scraper/backend/internal/models"
)

// rssSource defines a single RSS feed to poll
type rssSource struct {
	URL  string
	Team models.Team
	Name string
}

var rssSources = []rssSource{
	{
		URL:  "https://www.mlb.com/orioles/rss/news",
		Team: models.TeamOrioles,
		Name: "MLB Orioles",
	},
	{
		URL:  "https://www.nfl.com/rss/rsslanding?searchString=team/ravens",
		Team: models.TeamRavens,
		Name: "NFL Ravens",
	},
	{
		URL:  "https://www.espn.com/espn/rss/mlb/news",
		Team: models.TeamOrioles,
		Name: "ESPN MLB",
	},
	{
		URL:  "https://www.espn.com/espn/rss/nfl/news",
		Team: models.TeamRavens,
		Name: "ESPN NFL",
	},
	{
		URL:  "https://feeds.baltimoresun.com/baltimoresun/sports/ravens",
		Team: models.TeamRavens,
		Name: "Baltimore Sun Ravens",
	},
	{
		URL:  "https://feeds.baltimoresun.com/baltimoresun/sports/orioles",
		Team: models.TeamOrioles,
		Name: "Baltimore Sun Orioles",
	},
}

// RSSKeywords for determining team if not explicit
var oriolesKeywords = []string{"orioles", "o's", "baltimore baseball", "camden yards"}
var ravensKeywords = []string{"ravens", "lamar", "harbaugh", "m&t bank"}

type RSSScraper struct {
	parser *gofeed.Parser
}

func NewRSSScraper() *RSSScraper {
	return &RSSScraper{parser: gofeed.NewParser()}
}

func (s *RSSScraper) Name() string            { return "RSS" }
func (s *RSSScraper) Interval() time.Duration { return 10 * time.Minute }

func (s *RSSScraper) Fetch(ctx context.Context) ([]models.Article, error) {
	var articles []models.Article

	// Each feed is fetched concurrently
	type result struct {
		items []models.Article
		err   error
	}
	ch := make(chan result, len(rssSources))

	for _, src := range rssSources {
		src := src
		go func() {
			items, err := s.fetchFeed(ctx, src)
			ch <- result{items, err}
		}()
	}

	for range rssSources {
		res := <-ch
		if res.err == nil {
			articles = append(articles, res.items...)
		}
	}

	return articles, nil
}

func (s *RSSScraper) fetchFeed(ctx context.Context, src rssSource) ([]models.Article, error) {
	feed, err := s.parser.ParseURLWithContext(src.URL, ctx)
	if err != nil {
		return nil, err
	}

	var articles []models.Article
	for _, item := range feed.Items {
		pub := time.Now()
		if item.PublishedParsed != nil {
			pub = *item.PublishedParsed
		}

		// Use feed's team, or try to detect from title
		team := src.Team
		if team == "" {
			team = detectTeam(item.Title + " " + item.Description)
		}
		if team == "" {
			continue // Skip if we can't tie to Orioles or Ravens
		}

		imgURL := ""
		if item.Image != nil {
			imgURL = item.Image.URL
		}

		articles = append(articles, models.Article{
			Title:       item.Title,
			URL:         item.Link,
			Summary:     item.Description,
			Source:      models.SourceRSS,
			Team:        team,
			Author:      strings.Join(item.Authors.Names(), ", "),
			ImageURL:    imgURL,
			PublishedAt: pub,
		})
	}
	return articles, nil
}

func detectTeam(text string) models.Team {
	lower := strings.ToLower(text)
	isOrioles := containsAny(lower, oriolesKeywords)
	isRavens := containsAny(lower, ravensKeywords)
	if isOrioles && isRavens {
		return models.TeamBoth
	}
	if isOrioles {
		return models.TeamOrioles
	}
	if isRavens {
		return models.TeamRavens
	}
	return ""
}

func containsAny(s string, keywords []string) bool {
	for _, kw := range keywords {
		if strings.Contains(s, kw) {
			return true
		}
	}
	return false
}

// Helper to get author names from gofeed persons slice
func (p gofeed.Persons) Names() []string {
	names := make([]string, 0, len(p))
	for _, person := range p {
		if person.Name != "" {
			names = append(names, person.Name)
		}
	}
	return names
}
