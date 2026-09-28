import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";
import { FullFilter, constructFilterSort } from "./types";

export const load: PageLoad = async ({ parent, params, url }) => {
	const data = await parent();

	// /history/tracks resolves the user from the auth token rather than the
	// path (apis/history.go), so this page can only ever show your own
	// history. Refuse to render it under someone else's profile.
	if (!data.user || params.id !== data.user.id) {
		throw error(404, { message: "Not found" });
	}

	const query: Record<string, string> = {};

	const filter = FullFilter.parse({
		query: url.searchParams.get("query") ?? "",
		sort: url.searchParams.get("sort") ?? undefined,
		status: url.searchParams.get("status") ?? undefined,
	});

	constructFilterSort(filter, query);

	const history = await data.apiClient.getTrackHistory({ query });
	if (!history.success) {
		throw error(history.error.code, { message: history.error.message });
	}

	return {
		...data,
		page: history.data.page,
		history: history.data.history,
		filter,
	};
};
