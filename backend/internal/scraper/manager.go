package scraper

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	goredis "github.com/redis/go-redis/v9"

	"github.com/kaushikduvur/baltimore-sports-scraper/backend/internal/models"
	redisClient "github.com/kaushikduvur/baltimore-sports-scraper/backend/internal/redis"
)

// Scraper is the interface every data source must implement
type Scraper interface {
	Name() string
	Fetch(ctx context.Context) ([]models.Article, error)
	Interval() time.Duration
}

// Manager runs all scrapers concurrently using goroutines
type Manager struct {
	pool     *pgxpool.Pool
	rdb      *goredis.Client
	scrapers []Scraper
}

func NewManager(pool *pgxpool.Pool, rdb *goredis.Client) *Manager {
	return &Manager{
		pool: pool,
		rdb:  rdb,
		scrapers: []Scraper{
			NewRSSScraper(),
			NewGDELTScraper(),
			// NewRedditScraper(),   // Uncomment when Reddit credentials are set
			// NewBlueskyScraper(),  // Uncomment when Bluesky credentials are set
		},
	}
}

// Start launches each scraper in its own goroutine
func (m *Manager) Start(ctx context.Context) {
	log.Printf("Starting %d scrapers", len(m.scrapers))
	for _, s := range m.scrapers {
		go m.runScraper(ctx, s)
	}
}

func (m *Manager) runScraper(ctx context.Context, s Scraper) {
	log.Printf("[%s] Scraper started, interval: %v", s.Name(), s.Interval())
	ticker := time.NewTicker(s.Interval())
	defer ticker.Stop()

	// Run immediately on start
	m.executeAndStore(ctx, s)

	for {
		select {
		case <-ctx.Done():
			log.Printf("[%s] Scraper stopped", s.Name())
			return
		case <-ticker.C:
			m.executeAndStore(ctx, s)
		}
	}
}

func (m *Manager) executeAndStore(ctx context.Context, s Scraper) {
	articles, err := s.Fetch(ctx)
	if err != nil {
		log.Printf("[%s] Fetch error: %v", s.Name(), err)
		return
	}

	newCount := 0
	for _, article := range articles {
		inserted, err := m.upsertArticle(ctx, article)
		if err != nil {
			log.Printf("[%s] DB upsert error: %v", s.Name(), err)
			continue
		}
		if inserted {
			newCount++
			// Publish to Redis Stream so SSE clients get notified
			if err := redisClient.PublishArticle(ctx, m.rdb, article); err != nil {
				log.Printf("[%s] Redis publish error: %v", s.Name(), err)
			}
		}
	}

	if newCount > 0 {
		log.Printf("[%s] Stored %d new articles", s.Name(), newCount)
	}
}

// upsertArticle inserts an article, ignoring duplicates (by URL). Returns true if newly inserted.
func (m *Manager) upsertArticle(ctx context.Context, a models.Article) (bool, error) {
	tag, err := m.pool.Exec(ctx, `
		INSERT INTO articles (title, url, summary, source, team, author, image_url, sentiment_score, published_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (url) DO NOTHING
	`, a.Title, a.URL, a.Summary, a.Source, a.Team, a.Author, a.ImageURL, a.SentimentScore, a.PublishedAt)

	if err != nil {
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}
