import type { PageServerLoad } from './$types';

const API_URL = process.env.API_URL ?? 'http://localhost:8080';

export const load: PageServerLoad = async ({ fetch, url }) => {
	const team = url.searchParams.get('team') ?? '';
	const limit = 40;

	const params = new URLSearchParams({ limit: String(limit) });
	if (team) params.set('team', team);

	try {
		const res = await fetch(`${API_URL}/api/articles?${params}`);
		if (!res.ok) throw new Error(`API error: ${res.status}`);

		const data = await res.json();
		return {
			articles: data.articles ?? [],
			team,
			error: null
		};
	} catch (err) {
		console.error('Failed to load articles:', err);
		return {
			articles: [],
			team,
			error: 'Failed to load articles. Is the backend running?'
		};
	}
};
