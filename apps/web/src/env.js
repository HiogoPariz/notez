import { defineEnvVars } from '@sveltejs/kit/env';

export const variables = defineEnvVars({
	// Lida em runtime; ausente cai para o fallback do proxy.
	API_URL: { schema: (input) => input ?? '' },
	// Inlinadas em build-time (server-only).
	GOOGLE_CLIENT_ID: { static: true },
	GOOGLE_CLIENT_SECRET: { static: true },
	// Consumidas pelo @auth/sveltekit (via $env/dynamic/private).
	AUTH_SECRET: {},
	AUTH_TRUST_HOST: { schema: (input) => input ?? 'false' }
});
