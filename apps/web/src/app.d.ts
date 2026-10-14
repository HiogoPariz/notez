// See https://kit.svelte.dev/docs/types#app
import type { Session } from '@auth/core/types';

// `@auth/sveltekit` already augments `App.Locals` with `auth()`, `signIn` and
// `signOut` (and `App.PageData.session`). Here we only add the `session` field
// that `src/hooks.server.js` populates after a successful authentication.
declare global {
	namespace App {
		// interface Error {}
		interface Locals {
			session?: Session | null;
		}
		// interface PageData {}
		// interface PageState {}
		// interface Platform {}
	}
}

export {};
