# Baltimore Sports Hub 

A real-time Baltimore Orioles & Ravens news aggregator built with a modern, hireable stack. Articles are scraped from multiple sources concurrently, stored in PostgreSQL, and pushed live to connected browsers via Server-Sent Events — no refresh needed.

---

## Stack

| Layer | Technology | Why |
|---|---|---|
| Frontend | SvelteKit (SSR) | SSR for SEO + instant hydration |
| Backend | Go | Goroutines make concurrent scraping trivial |
| Real-time | Redis Streams + SSE | New articles pushed to browsers instantly |
| Database | PostgreSQL + Supabase | Reliable, structured, real-time-ready |
| Scraping | Go HTTP + gofeed | Fast, native, zero overhead |
| Social | Bluesky AT Protocol + Reddit | Both fully free APIs |
| News | GDELT Project API | Academic-grade, thousands of outlets, no key required |
| Hosting | Fly.io (backend) + Vercel (frontend) | Both have generous free tiers |

---

## Data Sources

| Source | What you get | Status |
|---|---|---|
| **GDELT API** | Thousands of news articles with sentiment scores | ✅ Active |
| **ESPN RSS** | Official game news, injury reports | ✅ Active |
| **MLB.com RSS** | Official Orioles news | ✅ Active |
| **NFL.com RSS** | Official Ravens news | ✅ Active |
| **Baltimore Sun RSS** | Local coverage for both teams | ✅ Active |
| **Reddit API** | r/orioles, r/ravens posts & comments | 🔧 Stub (needs credentials) |
| **Bluesky AT Protocol** | Public posts about Orioles/Ravens | 🔧 Stub (needs credentials) |
| **YouTube Data API v3** | Orioles/Ravens video content | 🔧 Planned |

---

## Project Structure

```
baltimore-sports-hub/
├── backend/                        # Go microservice
│   ├── cmd/server/main.go          # Entry point, router, graceful shutdown
│   ├── internal/
│   │   ├── api/
│   │   │   ├── handlers.go         # REST endpoints (articles, search)
│   │   │   └── sse.go              # Server-Sent Events stream
│   │   ├── db/postgres.go          # Connection pool + auto-migrations
│   │   ├── models/article.go       # Unified Article model
│   │   ├── redis/streams.go        # Publish/subscribe via Redis Streams
│   │   └── scraper/
│   │       ├── manager.go          # Concurrent scraper orchestrator
│   │       ├── rss.go              # ESPN, MLB, NFL, Baltimore Sun feeds
│   │       ├── gdelt.go            # GDELT DOC 2.0 API
│   │       ├── reddit.go           # Reddit API (stub)
│   │       └── bluesky.go          # Bluesky AT Protocol (stub)
│   ├── Dockerfile
│   └── go.mod
├── frontend/                       # SvelteKit app
│   ├── src/
│   │   ├── routes/
│   │   │   ├── +layout.svelte      # App shell, live indicator, SSE mount
│   │   │   ├── +page.server.ts     # SSR article load
│   │   │   └── +page.svelte        # Main feed page
│   │   └── lib/
│   │       ├── components/
│   │       │   ├── ArticleCard.svelte
│   │       │   ├── SearchBar.svelte
│   │       │   └── LiveFeed.svelte  # Headless SSE subscriber
│   │       └── stores/articles.ts  # Svelte stores + derived filters
│   ├── svelte.config.js
│   └── package.json
├── docker-compose.yml              # Local dev stack (Postgres + Redis + backend + frontend)
├── .env.example                    # All required environment variables
└── ROADMAP.md
```

---

## Quick Start

### Prerequisites
- [Docker & Docker Compose](https://docs.docker.com/get-docker/)
- [Go 1.22+](https://go.dev/dl/) (for local backend dev)
- [Node 20+](https://nodejs.org/) (for local frontend dev)

### 1. Clone & configure

```bash
git clone https://github.com/kaushikduvur/Baltimore_Sports_Scraper.git
cd Baltimore_Sports_Scraper
cp .env.example .env
```

### 2. Start everything with Docker

```bash
docker-compose up --build
```

This starts PostgreSQL, Redis, the Go backend, and the SvelteKit frontend.

- Frontend → http://localhost:5173
- Backend API → http://localhost:8080
- Health check → http://localhost:8080/health

### 3. Local development (without Docker)

**Backend:**
```bash
cd backend
go mod tidy
go run ./cmd/server
```

**Frontend:**
```bash
cd frontend
npm install
npm run dev
```

---

## API Endpoints

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/health` | Health check |
| `GET` | `/api/articles` | List articles (`?team=orioles&source=rss&limit=20&offset=0`) |
| `GET` | `/api/articles/:id` | Single article by ID |
| `GET` | `/api/search` | Full-text search (`?q=lamar+jackson&team=ravens`) |
| `GET` | `/api/stream` | SSE live stream (`?team=orioles`) |

---

## Deployment

**Backend → Fly.io:**
```bash
cd backend
flyctl launch
flyctl deploy
```

**Frontend → Vercel:**
```bash
cd frontend
vercel deploy
```

Set `VITE_API_URL` and `API_URL` in Vercel's environment variables to your Fly.io backend URL.

---

## Adding More Data Sources

1. Create a new file in `backend/internal/scraper/` implementing the `Scraper` interface:
   ```go
   type Scraper interface {
       Name() string
       Fetch(ctx context.Context) ([]models.Article, error)
       Interval() time.Duration
   }
   ```
2. Register it in `manager.go`'s `NewManager()` scrapers slice.
3. That's it — the manager handles concurrent scheduling, deduplication, and Redis publishing automatically.

---


