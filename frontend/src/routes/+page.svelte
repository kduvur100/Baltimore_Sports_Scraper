<script lang="ts">
	import { onMount } from 'svelte';
	import { articles, filteredArticles } from '$lib/stores/articles';
	import ArticleCard from '$lib/components/ArticleCard.svelte';
	import SearchBar from '$lib/components/SearchBar.svelte';
	import type { PageData } from './$types';

	export let data: PageData;

	// Seed the store with SSR-loaded articles
	onMount(() => {
		articles.set(data.articles ?? []);
	});
</script>

<svelte:head>
	<title>Baltimore Sports Hub — Live Orioles & Ravens News</title>
</svelte:head>

<div class="page">
	<div class="page-header">
		<h1 class="page-title">Baltimore Sports Hub</h1>
		<p class="page-subtitle">
			Real-time news from GDELT, ESPN, Baltimore Sun, MLB, NFL, Reddit & Bluesky
		</p>
	</div>

	<SearchBar />

	{#if data.error}
		<div class="error-banner">
			⚠️ {data.error}
		</div>
	{/if}

	{#if $filteredArticles.length === 0}
		<div class="empty-state">
			<p>No articles found. Check back soon or try a different filter.</p>
		</div>
	{:else}
		<div class="article-grid">
			{#each $filteredArticles as article (article.id)}
				<ArticleCard {article} />
			{/each}
		</div>
	{/if}
</div>

<style>
	.page {
		padding-bottom: 2rem;
	}
	.page-header {
		margin-bottom: 1.5rem;
	}
	.page-title {
		font-size: 1.75rem;
		font-weight: 800;
		letter-spacing: -0.02em;
		margin-bottom: 0.25rem;
	}
	.page-subtitle {
		color: #666;
		font-size: 0.9rem;
	}
	.error-banner {
		background: #fff3cd;
		border: 1px solid #ffc107;
		border-radius: 8px;
		padding: 0.75rem 1rem;
		margin-bottom: 1rem;
		font-size: 0.9rem;
	}
	.empty-state {
		text-align: center;
		padding: 3rem;
		color: #888;
	}
	.article-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
		gap: 1rem;
	}
</style>
