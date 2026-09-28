import type { Album, Page } from "$lib/api/types";
import type { PageLoad } from "./$types";

const SEARCH_ALBUMS_PAGE_SIZE = 30;

export const load: PageLoad = async ({ parent, url }) => {
	const data = await parent();

	const query = url.searchParams.get("query") ?? "";

	let albums = [] as Album[];
	let page: Page | null = null;

	if (query) {
		const res = await data.apiClient.searchAlbums({
			query: { query, perPage: SEARCH_ALBUMS_PAGE_SIZE.toString() },
		});

		if (res.success) {
			albums = res.data.albums;
			page = res.data.page;
		}
	}

	return {
		...data,
		query,
		albums,
		page,
	};
};
