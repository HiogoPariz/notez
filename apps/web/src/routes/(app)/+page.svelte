<script>
	// @ts-nocheck

	import { saveNote } from '$lib/services/saveNote';

	// @ts-nocheck
	import throttle from 'lodash.throttle';
	import { onMount } from 'svelte';
	/**
	 * @type {import("@editorjs/editorjs").default}
	 */
	let editor;
	let noteTitle;

	// AUTO SAVE FUNCTION
	const DELAY = 300;

	const handleSave = throttle(() => {
		save();
	}, DELAY);

	async function save() {
		try {
			const data = await editor.save();
			await saveNote({ title: noteTitle, content: JSON.stringify(data) });
		} catch (error) {
			console.error('ERROR:', error);
		}
	}

	onMount(async () => {
		const EditorJS = (await import('@editorjs/editorjs')).default;
		const Header = (await import('@editorjs/header')).default;
		const Checklist = (await import('@editorjs/checklist')).default;
		const List = (await import('@editorjs/list')).default;
		const Quote = (await import('@editorjs/quote')).default;
		editor = new EditorJS({
			autofocus: true,
			holder: 'editor',
			tools: {
				header: Header,
				checklist: {
					class: Checklist,
					inlineToolbar: true,
					shortcut: 'CMD+SHIFT+C'
				},
				list: {
					class: List,
					inlineToolbar: true,
					shortcut: 'CMD+SHIFT+L',
					config: {
						defaultStyle: 'unordered'
					}
				},
				quote: {
					class: Quote,
					inlineToolbar: true,
					shortcut: 'CMD+SHIFT+O'
				}
			},
			onChange: handleSave
		});
	});
</script>

<svelte:head>
	<title>{noteTitle || 'Default - Notez'}</title>
</svelte:head>
<section class="flex flex-col items-center p-4">
	<input bind:value={noteTitle} placeholder="Notez Title" />
	<div class="editor w-full" id="editor"></div>
</section>
