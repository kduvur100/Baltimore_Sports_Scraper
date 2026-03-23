import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [sveltekit()],
	server: {
		port: 5173,
		proxy: {
			// Proxy API calls to the Go backend in development
			'/api': {
				target: process.env.VITE_API_URL || 'http://localhost:8080',
				changeOrigin: true
			}
		}
	}
});
