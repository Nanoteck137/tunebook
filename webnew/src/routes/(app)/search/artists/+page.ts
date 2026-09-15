import type { Artist, Page } from "$lib/api/types";
import type { PageLoad } from "./$types";

const SEARCH_ARTISTS_PAGE_SIZE = 30;

export const load: PageLoad = async ({ parent, url }) => {
	const data = await parent();

	const query = url.searchParams.get("query") ?? "";

	let artists = [] as Artist[];
	let page: Page | null = null;

	if (query) {
		const res = await data.apiClient.searchArtists({
			query: { query, perPage: SEARCH_ARTISTS_PAGE_SIZE.toString() },
		});

		if (res.success) {
			artists = res.data.artists;
			page = res.data.page;
		}
	}

	return {
		...data,
		query,
		artists,
		page,
	};
};
