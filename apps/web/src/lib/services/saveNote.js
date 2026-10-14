/**
 * Envia uma nota para a API (via proxy BFF em /api/note).
 * @param {{ title: string, content: string }} note
 */
export async function saveNote(note) {
	const res = await fetch('/api/note', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify(note)
	});

	return await res.json();
}
