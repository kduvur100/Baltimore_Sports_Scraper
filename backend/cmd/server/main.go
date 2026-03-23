package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"github.com/kaushikduvur/baltimore-sports-scraper/backend/internal/api"
	"github.com/kaushikduvur/baltimore-sports-scraper/backend/internal/db"
	"github.com/kaushikduvur/baltimore-sports-scraper/backend/internal/redis"
	"github.com/kaushikduvur/baltimore-sports-scraper/backend/internal/scraper"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize DB
	pool, err := db.NewPool(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	if err := db.RunMigrations(pool); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize Redis
	rdb := redis.NewClient(os.Getenv("REDIS_URL"))
	defer rdb.Close()

	// Initialize scraper manager — runs all scrapers concurrently
	sm := scraper.NewManager(pool, rdb)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go sm.Start(ctx)

	// Initialize HTTP router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{os.Getenv("FRONTEND_URL"), "http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		AllowCredentials: false,
	}))

	// Register API routes
	h := api.NewHandler(pool, rdb)
	r.Get("/health", h.Health)
	r.Get("/api/articles", h.GetArticles)
	r.Get("/api/articles/{id}", h.GetArticle)
	r.Get("/api/search", h.Search)
	r.Get("/api/stream", h.SSEStream) // Server-Sent Events

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 0, // SSE connections stay open
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited cleanly")
}
