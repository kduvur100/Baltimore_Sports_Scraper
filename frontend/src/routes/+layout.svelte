<script lang="ts">
	import LiveFeed from '$lib/components/LiveFeed.svelte';
	import { isLive, newArticleCount, clearNewCount } from '$lib/stores/articles';
</script>

<LiveFeed />

<div class="app-shell">
	<header class="app-header">
		<div class="header-inner">
			<a href="/" class="logo">
				<span class="logo-icon">⚾🏈</span>
				<span class="logo-text">Baltimore Sports Hub</span>
			</a>

			<div class="header-right">
				{#if $newArticleCount > 0}
					<button class="new-badge" on:click={clearNewCount}>
						+{$newArticleCount} new
					</button>
				{/if}
				<div class="live-indicator" class:live={$isLive}>
					<span class="live-dot"></span>
					{$isLive ? 'Live' : 'Connecting…'}
				</div>
			</div>
		</div>
	</header>

	<main class="main-content">
		<slot />
	</main>

	<footer class="app-footer">
		<p>Data sourced from GDELT, ESPN, Baltimore Sun, MLB, NFL, Reddit & Bluesky · Built with Go + SvelteKit</p>
	</footer>
</div>

<style>
	:global(*, *::before, *::after) { box-sizing: border-box; margin: 0; padding: 0; }
	:global(body) {
		font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
		background: #f4f5f7;
		color: #111;
	}

	.app-shell {
		min-height: 100vh;
		display: flex;
		flex-direction: column;
	}
	.app-header {
		background: #111;
		color: #fff;
		position: sticky;
		top: 0;
		z-index: 100;
		box-shadow: 0 2px 8px rgba(0,0,0,0.2);
	}
	.header-inner {
		max-width: 1200px;
		margin: 0 auto;
		padding: 0.75rem 1.5rem;
		display: flex;
		justify-content: space-between;
		align-items: center;
	}
	.logo {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		text-decoration: none;
		color: #fff;
	}
	.logo-icon { font-size: 1.4rem; }
	.logo-text { font-size: 1.1rem; font-weight: 700; letter-spacing: -0.01em; }

	.header-right {
		display: flex;
		align-items: center;
		gap: 0.75rem;
	}
	.new-badge {
		background: #DF4601;
		color: #fff;
		border: none;
		border-radius: 999px;
		padding: 0.3rem 0.8rem;
		font-size: 0.8rem;
		font-weight: 600;
		cursor: pointer;
		animation: pulse 1.5s infinite;
	}
	@keyframes pulse {
		0%, 100% { opacity: 1; }
		50% { opacity: 0.7; }
	}
	.live-indicator {
		display: flex;
		align-items: center;
		gap: 0.4rem;
		font-size: 0.8rem;
		color: #aaa;
	}
	.live-indicator.live { color: #4ade80; }
	.live-dot {
		width: 8px; height: 8px;
		border-radius: 50%;
		background: #666;
	}
	.live-indicator.live .live-dot {
		background: #4ade80;
		box-shadow: 0 0 0 2px rgba(74,222,128,0.3);
		animation: blink 1.5s infinite;
	}
	@keyframes blink {
		0%, 100% { opacity: 1; }
		50% { opacity: 0.4; }
	}

	.main-content {
		flex: 1;
		max-width: 1200px;
		margin: 0 auto;
		padding: 1.5rem;
		width: 100%;
	}
	.app-footer {
		text-align: center;
		padding: 1.5rem;
		font-size: 0.8rem;
		color: #999;
		border-top: 1px solid #e5e5e5;
	}
</style>
