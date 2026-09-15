import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";
import { FullFilter, constructFilterSort } from "./types";

export const load: PageLoad = async ({ parent, url }) => {
	const data = await parent();

	const query: Record<string, string> = {};

	const filter = FullFilter.parse({
		query: url.searchParams.get("query") ?? "",
		sort: url.searchParams.get("sort") ?? undefined,
		filters: {
			tags: url.searchParams.get("tags")?.split(",") ?? [],
		},
		excludes: {
			tags: url.searchParams.get("excludeTags")?.split(",") ?? [],
		},
	});

	constructFilterSort(filter, query);

	const artists = await data.apiClient.getArtists({
		query,
	});
	if (!artists.success) {
		throw error(artists.error.code, artists.error.message);
	}

	return {
		...data,
		page: artists.data.page,
		artists: artists.data.artists,
		filter,
	};
};
