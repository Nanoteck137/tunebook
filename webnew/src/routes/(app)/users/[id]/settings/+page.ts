import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, params }) => {
	const data = await parent();

	// Every mutation on this page targets the caller's own account, so the tab
	// is hidden for other profiles (users/[id]/+layout.svelte). The route itself
	// was still reachable directly, which rendered your own tokens and settings
	// under someone else's profile header. Refuse to render it.
	if (!data.user || params.id !== data.user.id) {
		throw error(404, { message: "Not found" });
	}

	const res = await data.apiClient.getApiTokens();
	if (!res.success) {
		throw error(res.error.code, { message: res.error.message });
	}

	return {
		...data,

		tokens: res.data.tokens,
	};
};
