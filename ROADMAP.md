# Baltimore Sports Hub — Roadmap

This document tracks what's built, what's next, and future ideas.
Update this as you ship features.

---

## ✅ Day 1 — Foundation (Done)

- [x] Go backend with chi router, graceful shutdown
- [x] PostgreSQL schema with auto-migrations
- [x] Redis Streams event bus
- [x] Concurrent scraper manager (goroutines)
- [x] RSS scraper — ESPN, MLB, NFL, Baltimore Sun (live)
- [x] GDELT scraper — global news with sentiment scores (live)
- [x] Server-Sent Events (SSE) endpoint for real-time push
- [x] SvelteKit SSR frontend — SEO-ready, instant load
- [x] Live feed indicator + new article badge
- [x] Team filter (Orioles / Ravens / All)
- [x] Client-side search
- [x] ArticleCard with sentiment badge, team color accent
- [x] Docker Compose local dev stack
- [x] .env.example with all config documented
- [x] README with quick start + API docs

---

## 🔧 Day 2 — Social Sources

- [ ] **Reddit scraper** — activate `reddit.go` stub
  - Fetch top posts from r/orioles and r/ravens
  - Use PRAW-style Go client (no library needed, direct OAuth)
  - Store post text as article summary, link to thread
- [ ] **Bluesky AT Protocol scraper** — activate `bluesky.go` stub
  - `app.bsky.feed.searchPosts` to find #Orioles and #Ravens posts
  - Only surface posts with engagement (likes > 5) to reduce noise
- [ ] **Deduplication improvements**
  - Fuzzy title matching to catch same story from multiple sources
  - Use pg_trgm extension for similarity queries

---

## 🔍 Day 3 — Search Upgrade (Meilisearch)

- [ ] Add Meilisearch to Docker Compose
- [ ] Index all articles into Meilisearch on insert
- [ ] Replace PostgreSQL ILIKE search with Meilisearch instant search
- [ ] Add typo tolerance, faceted search (by team, source, date)
- [ ] Add search result highlighting in ArticleCard

---

## 📺 Day 4 — YouTube Integration

- [ ] YouTube Data API v3 scraper
  - Search for official Orioles/Ravens channels
  - Pull latest video uploads + highlights
  - Embed YouTube thumbnails in ArticleCard
- [ ] Video card variant in the grid (different layout from text articles)

---

## 🎨 Day 5 — UX Polish

- [ ] Infinite scroll / pagination
- [ ] Article skeleton loading states
- [ ] Dark mode toggle
- [ ] Mobile-responsive improvements
- [ ] Orioles / Ravens themed color scheme per tab
- [ ] Toasts for new articles (non-intrusive notification)
- [ ] Share button per article

---

## 📊 Day 6 — Analytics & Insights

- [ ] Sentiment trend chart (GDELT scores over time)
  - "Ravens fan mood this week" visualization
- [ ] Most mentioned players (NLP keyword extraction in Go)
- [ ] Source breakdown pie chart
- [ ] Article volume per day heatmap

---

## 🚀 Day 7 — Production Deployment

- [ ] Backend → Fly.io
  - `fly launch`, set secrets, deploy
  - Fly Postgres or connect to Supabase
  - Fly Redis (Upstash)
- [ ] Frontend → Vercel
  - Set `VITE_API_URL` env var
  - Enable Vercel Analytics
- [ ] Custom domain
- [ ] Set up GitHub Actions CI
  - `go test ./...` on every push
  - `svelte-check` on every push
  - Docker build validation

---

## 💡 Future Ideas

- **Email digest** — daily email with top 5 stories (SendGrid + cron)
- **Push notifications** — browser push for breaking news
- **Score widget** — live game scores via ESPN API
- **Discord/Slack bot** — post breaking news to a channel
- **Comment aggregation** — show Reddit comment count on each article
- **Meilisearch semantic search** — vector search for "articles like this"
- **Playwright scraper** — for sites without RSS (The Athletic paywall bypass via headlines only)
