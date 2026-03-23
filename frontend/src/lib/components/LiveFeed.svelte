<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { isLive, prependArticle, teamFilter } from '$lib/stores/articles';
	import type { Article } from '$lib/stores/articles';

	const API_URL = import.meta.env.VITE_API_URL ?? '';

	let eventSource: EventSource | null = null;
	let retryTimeout: ReturnType<typeof setTimeout>;

	function connect() {
		const team = $teamFilter;
		const url = `${API_URL}/api/stream${team ? `?team=${team}` : ''}`;

		eventSource = new EventSource(url);

		eventSource.addEventListener('connected', () => {
			isLive.set(true);
			console.log('[SSE] Connected to live feed');
		});

		eventSource.addEventListener('article', (e: MessageEvent) => {
			try {
				const article: Article = JSON.parse(e.data);
				prependArticle(article);
			} catch (err) {
				console.error('[SSE] Failed to parse article:', err);
			}
		});

		eventSource.onerror = () => {
			isLive.set(false);
			eventSource?.close();
			// Reconnect after 5s
			retryTimeout = setTimeout(connect, 5000);
		};
	}

	onMount(connect);

	onDestroy(() => {
		eventSource?.close();
		clearTimeout(retryTimeout);
		isLive.set(false);
	});
</script>

<!-- This component has no UI — it's a headless SSE subscriber -->
<!-- The live status indicator is rendered in +layout.svelte -->
