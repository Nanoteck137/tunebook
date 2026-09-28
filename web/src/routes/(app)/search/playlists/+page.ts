import type { Page, Playlist } from "$lib/api/types";
import type { PageLoad } from "./$types";

const SEARCH_PLAYLISTS_PAGE_SIZE = 30;

export const load: PageLoad = async ({ parent, url }) => {
	const data = await parent();

	const query = url.searchParams.get("query") ?? "";

	let playlists = [] as Playlist[];
	let page: Page | null = null;

	if (query) {
		const res = await data.apiClient.searchPlaylists({
			query: { query, perPage: SEARCH_PLAYLISTS_PAGE_SIZE.toString() },
		});

		if (res.success) {
			playlists = res.data.playlists;
			page = res.data.page;
		}
	}

	return {
		...data,
		query,
		playlists,
		page,
	};
};
