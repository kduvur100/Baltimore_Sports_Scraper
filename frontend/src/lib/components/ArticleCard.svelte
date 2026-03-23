<script lang="ts">
	import { formatDistanceToNow } from 'date-fns';
	import type { Article } from '$lib/stores/articles';

	export let article: Article;

	// Source → label map
	const sourceLabels: Record<string, string> = {
		rss: 'RSS',
		gdelt: 'GDELT',
		reddit: 'Reddit',
		bluesky: 'Bluesky',
		youtube: 'YouTube'
	};

	// Team → color accent
	const teamColors: Record<string, string> = {
		orioles: '#DF4601', // Orioles orange
		ravens: '#241773',  // Ravens purple
		both: '#888'
	};

	$: timeAgo = formatDistanceToNow(new Date(article.published_at), { addSuffix: true });
	$: accent = teamColors[article.team] ?? '#888';
	$: sourceLabel = sourceLabels[article.source] ?? article.source;

	// Sentiment: 0–1 → negative/neutral/positive label
	$: sentimentLabel =
		article.sentiment_score == null
			? null
			: article.sentiment_score > 0.55
			? '😊 Positive'
			: article.sentiment_score < 0.45
			? '😟 Negative'
			: '😐 Neutral';
</script>

<article
	class="card"
	style="--accent: {accent};"
>
	{#if article.image_url}
		<a href={article.url} target="_blank" rel="noopener noreferrer" class="card-image-link">
			<img
				src={article.image_url}
				alt={article.title}
				class="card-image"
				loading="lazy"
			/>
		</a>
	{/if}

	<div class="card-body">
		<div class="card-meta">
			<span class="badge team-badge">{article.team.toUpperCase()}</span>
			<span class="badge source-badge">{sourceLabel}</span>
			{#if sentimentLabel}
				<span class="badge sentiment-badge">{sentimentLabel}</span>
			{/if}
		</div>

		<h2 class="card-title">
			<a href={article.url} target="_blank" rel="noopener noreferrer">
				{article.title}
			</a>
		</h2>

		{#if article.summary}
			<p class="card-summary">{article.summary}</p>
		{/if}

		<footer class="card-footer">
			{#if article.author}
				<span class="author">by {article.author}</span>
			{/if}
			<time datetime={article.published_at}>{timeAgo}</time>
		</footer>
	</div>
</article>

<style>
	.card {
		border-left: 4px solid var(--accent);
		background: #fff;
		border-radius: 8px;
		box-shadow: 0 1px 4px rgba(0,0,0,0.08);
		overflow: hidden;
		transition: box-shadow 0.2s;
	}
	.card:hover {
		box-shadow: 0 4px 16px rgba(0,0,0,0.13);
	}
	.card-image-link { display: block; }
	.card-image {
		width: 100%;
		height: 180px;
		object-fit: cover;
	}
	.card-body {
		padding: 1rem;
	}
	.card-meta {
		display: flex;
		gap: 0.4rem;
		flex-wrap: wrap;
		margin-bottom: 0.5rem;
	}
	.badge {
		font-size: 0.7rem;
		font-weight: 600;
		padding: 2px 8px;
		border-radius: 999px;
		text-transform: uppercase;
		letter-spacing: 0.04em;
	}
	.team-badge {
		background: var(--accent);
		color: #fff;
	}
	.source-badge {
		background: #f0f0f0;
		color: #555;
	}
	.sentiment-badge {
		background: #f5f5f5;
		color: #333;
	}
	.card-title {
		font-size: 1rem;
		font-weight: 600;
		margin: 0 0 0.5rem;
		line-height: 1.35;
	}
	.card-title a {
		color: #111;
		text-decoration: none;
	}
	.card-title a:hover {
		text-decoration: underline;
	}
	.card-summary {
		font-size: 0.875rem;
		color: #555;
		margin: 0 0 0.75rem;
		display: -webkit-box;
		-webkit-line-clamp: 3;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}
	.card-footer {
		display: flex;
		justify-content: space-between;
		font-size: 0.78rem;
		color: #888;
	}
	.author { font-style: italic; }
</style>
