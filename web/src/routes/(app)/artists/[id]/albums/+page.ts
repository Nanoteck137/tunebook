import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";
import { AlbumFilter, applySort } from "./types";

export const load: PageLoad = async ({ parent, params, url }) => {
	const data = await parent();

	const query: Record<string, string> = {};

	const filter = AlbumFilter.parse({
		sort: url.searchParams.get("sort") ?? undefined,
	});
	applySort(filter, query);

	const albums = await data.apiClient.getAlbums({
		query: {
			...query,
			filter: `artistId = "${params.id}" or featuringArtists has "${params.id}"`,
		},
	});
	if (!albums.success) {
		throw error(albums.error.code, { message: albums.error.message });
	}

	return {
		...data,
		page: albums.data.page,
		albums: albums.data.albums,
		filter,
	};
};
