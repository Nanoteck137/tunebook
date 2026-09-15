import type { Page, Track } from "$lib/api/types";
import type { PageLoad } from "./$types";

const SEARCH_TRACKS_PAGE_SIZE = 50;

export const load: PageLoad = async ({ parent, url }) => {
	const data = await parent();

	const query = url.searchParams.get("query") ?? "";

	let tracks = [] as Track[];
	let page: Page | null = null;

	if (query) {
		const res = await data.apiClient.searchTracks({
			query: { query, perPage: SEARCH_TRACKS_PAGE_SIZE.toString() },
		});

		if (res.success) {
			tracks = res.data.tracks;
			page = res.data.page;
		}
	}

	return {
		...data,
		query,
		tracks,
		page,
	};
};
