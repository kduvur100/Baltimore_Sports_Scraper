import { writable, derived } from 'svelte/store';

export type Team = 'orioles' | 'ravens' | 'both';
export type Source = 'rss' | 'gdelt' | 'reddit' | 'bluesky' | 'youtube';

export interface Article {
	id: string;
	title: string;
	url: string;
	summary: string;
	source: Source;
	team: Team;
	author?: string;
	image_url?: string;
	sentiment_score?: number;
	published_at: string;
	created_at: string;
}

// All articles loaded from the server
export const articles = writable<Article[]>([]);

// Active team filter ('orioles' | 'ravens' | '' = all)
export const teamFilter = writable<Team | ''>('');

// Search query string
export const searchQuery = writable('');

// Whether a live SSE connection is active
export const isLive = writable(false);

// Count of new articles received via SSE since last page load
export const newArticleCount = writable(0);

// Derived: articles filtered by team and search query
export const filteredArticles = derived(
	[articles, teamFilter, searchQuery],
	([$articles, $teamFilter, $searchQuery]) => {
		let result = $articles;

		if ($teamFilter) {
			result = result.filter(
				(a) => a.team === $teamFilter || a.team === 'both'
			);
		}

		if ($searchQuery.trim()) {
			const q = $searchQuery.toLowerCase();
			result = result.filter(
				(a) =>
					a.title.toLowerCase().includes(q) ||
					a.summary?.toLowerCase().includes(q)
			);
		}

		return result;
	}
);

// Add a new article to the top of the list (called when SSE event arrives)
export function prependArticle(article: Article) {
	articles.update((current) => {
		// Avoid duplicates
		if (current.find((a) => a.id === article.id)) return current;
		newArticleCount.update((n) => n + 1);
		return [article, ...current];
	});
}

// Reset the new article badge count
export function clearNewCount() {
	newArticleCount.set(0);
}
