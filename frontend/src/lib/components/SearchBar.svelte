<script lang="ts">
	import { searchQuery, teamFilter, clearNewCount } from '$lib/stores/articles';
	import type { Team } from '$lib/stores/articles';

	function handleTeam(team: Team | '') {
		teamFilter.set(team);
		clearNewCount();
	}
</script>

<div class="search-bar">
	<div class="search-input-wrap">
		<svg class="search-icon" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
			<circle cx="11" cy="11" r="8"/><path d="m21 21-4.35-4.35"/>
		</svg>
		<input
			type="search"
			placeholder="Search articles..."
			bind:value={$searchQuery}
			class="search-input"
		/>
	</div>

	<div class="team-tabs" role="group" aria-label="Filter by team">
		<button
			class="tab"
			class:active={$teamFilter === ''}
			on:click={() => handleTeam('')}
		>All</button>
		<button
			class="tab orioles"
			class:active={$teamFilter === 'orioles'}
			on:click={() => handleTeam('orioles')}
		>⚾ Orioles</button>
		<button
			class="tab ravens"
			class:active={$teamFilter === 'ravens'}
			on:click={() => handleTeam('ravens')}
		>🏈 Ravens</button>
	</div>
</div>

<style>
	.search-bar {
		display: flex;
		flex-wrap: wrap;
		gap: 0.75rem;
		align-items: center;
		margin-bottom: 1.25rem;
	}
	.search-input-wrap {
		position: relative;
		flex: 1;
		min-width: 200px;
	}
	.search-icon {
		position: absolute;
		left: 0.75rem;
		top: 50%;
		transform: translateY(-50%);
		color: #999;
		pointer-events: none;
	}
	.search-input {
		width: 100%;
		padding: 0.55rem 0.75rem 0.55rem 2.25rem;
		border: 1px solid #ddd;
		border-radius: 8px;
		font-size: 0.9rem;
		outline: none;
		transition: border-color 0.2s;
		box-sizing: border-box;
	}
	.search-input:focus { border-color: #555; }

	.team-tabs {
		display: flex;
		gap: 0.4rem;
	}
	.tab {
		padding: 0.45rem 1rem;
		border-radius: 999px;
		border: 1px solid #ddd;
		background: #f8f8f8;
		font-size: 0.85rem;
		cursor: pointer;
		font-weight: 500;
		transition: all 0.15s;
	}
	.tab:hover { background: #eee; }
	.tab.active { background: #111; color: #fff; border-color: #111; }
	.tab.orioles.active { background: #DF4601; border-color: #DF4601; }
	.tab.ravens.active  { background: #241773; border-color: #241773; }
</style>
