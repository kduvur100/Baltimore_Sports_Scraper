package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	goredis "github.com/redis/go-redis/v9"

	"github.com/kaushikduvur/baltimore-sports-scraper/backend/internal/models"
)

type Handler struct {
	db  *pgxpool.Pool
	rdb *goredis.Client
}

func NewHandler(db *pgxpool.Pool, rdb *goredis.Client) *Handler {
	return &Handler{db: db, rdb: rdb}
}

// Health godoc
// GET /health
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}

// GetArticles godoc
// GET /api/articles?team=orioles&source=rss&limit=20&offset=0
func (h *Handler) GetArticles(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	team := q.Get("team")     // "orioles", "ravens", "both", or "" (all)
	source := q.Get("source") // "rss", "gdelt", "reddit", "bluesky", or "" (all)
	limit := queryInt(q.Get("limit"), 20)
	offset := queryInt(q.Get("offset"), 0)

	if limit > 100 {
		limit = 100
	}

	// Build dynamic query
	args := []interface{}{}
	where := "WHERE 1=1"
	argIdx := 1

	if team != "" {
		where += " AND (team = $" + strconv.Itoa(argIdx) + " OR team = 'both')"
		args = append(args, team)
		argIdx++
	}
	if source != "" {
		where += " AND source = $" + strconv.Itoa(argIdx)
		args = append(args, source)
		argIdx++
	}

	args = append(args, limit, offset)
	query := `
		SELECT id, title, url, summary, source, team, author, image_url, sentiment_score, published_at, created_at
		FROM articles
		` + where + `
		ORDER BY published_at DESC
		LIMIT $` + strconv.Itoa(argIdx) + ` OFFSET $` + strconv.Itoa(argIdx+1)

	rows, err := h.db.Query(r.Context(), query, args...)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to fetch articles")
		return
	}
	defer rows.Close()

	articles := []models.Article{}
	for rows.Next() {
		var a models.Article
		if err := rows.Scan(
			&a.ID, &a.Title, &a.URL, &a.Summary, &a.Source, &a.Team,
			&a.Author, &a.ImageURL, &a.SentimentScore, &a.PublishedAt, &a.CreatedAt,
		); err != nil {
			continue
		}
		articles = append(articles, a)
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"articles": articles,
		"count":    len(articles),
		"offset":   offset,
	})
}

// GetArticle godoc
// GET /api/articles/{id}
func (h *Handler) GetArticle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var a models.Article
	err := h.db.QueryRow(r.Context(), `
		SELECT id, title, url, summary, source, team, author, image_url, sentiment_score, published_at, created_at
		FROM articles WHERE id = $1
	`, id).Scan(
		&a.ID, &a.Title, &a.URL, &a.Summary, &a.Source, &a.Team,
		&a.Author, &a.ImageURL, &a.SentimentScore, &a.PublishedAt, &a.CreatedAt,
	)
	if err != nil {
		writeError(w, http.StatusNotFound, "Article not found")
		return
	}

	writeJSON(w, http.StatusOK, a)
}

// Search godoc
// GET /api/search?q=lamar+jackson&team=ravens
func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	team := r.URL.Query().Get("team")
	limit := queryInt(r.URL.Query().Get("limit"), 20)

	if query == "" {
		writeError(w, http.StatusBadRequest, "q parameter is required")
		return
	}

	args := []interface{}{"%" + query + "%", "%" + query + "%", limit}
	whereTeam := ""
	if team != "" {
		whereTeam = " AND (team = $4 OR team = 'both')"
		args = append(args, team)
	}

	rows, err := h.db.Query(r.Context(), `
		SELECT id, title, url, summary, source, team, author, image_url, sentiment_score, published_at, created_at
		FROM articles
		WHERE (title ILIKE $1 OR summary ILIKE $2)`+whereTeam+`
		ORDER BY published_at DESC
		LIMIT $3
	`, args...)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Search failed")
		return
	}
	defer rows.Close()

	articles := []models.Article{}
	for rows.Next() {
		var a models.Article
		if err := rows.Scan(
			&a.ID, &a.Title, &a.URL, &a.Summary, &a.Source, &a.Team,
			&a.Author, &a.ImageURL, &a.SentimentScore, &a.PublishedAt, &a.CreatedAt,
		); err != nil {
			continue
		}
		articles = append(articles, a)
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"articles": articles,
		"query":    query,
		"count":    len(articles),
	})
}

// --- Helpers ---

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func queryInt(s string, def int) int {
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil || v < 0 {
		return def
	}
	return v
}
