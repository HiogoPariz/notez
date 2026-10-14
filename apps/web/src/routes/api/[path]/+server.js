import { API_URL as API_URL_ENV } from '$app/env/private';

const API_URL = API_URL_ENV || 'http://localhost:3000';

/**
 * Proxy catch-all que encaminha as requisições de /api/* para a API,
 * devolvendo a resposta real do backend para o cliente.
 * @param {import('@sveltejs/kit').RequestEvent} event
 */
export async function fallback({ request, params }) {
	const path = params.path;

	const headers = new Headers(request.headers);
	headers.set('Accept', 'application/json');
	headers.delete('host');
	headers.delete('content-length');

	const hasBody = request.method !== 'GET' && request.method !== 'HEAD';

	const upstream = await fetch(`${API_URL}/${path}`, {
		method: request.method,
		headers,
		body: hasBody ? await request.arrayBuffer() : undefined
	});

	const body = await upstream.arrayBuffer();

	return new Response(body, {
		status: upstream.status,
		headers: {
			'content-type': upstream.headers.get('content-type') ?? 'application/json'
		}
	});
}
