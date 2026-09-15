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
			decade: url.searchParams.get("decade") ?? undefined,
			tags: url.searchParams.get("tags")?.split(",") ?? [],
		},
		excludes: {
			tags: url.searchParams.get("excludeTags")?.split(",") ?? [],
		},
	});

	constructFilterSort(filter, query);

	const res = await data.apiClient.getAlbums({ query });
	if (!res.success) {
		throw error(res.error.code, { message: res.error.message });
	}

	return {
		...data,
		page: res.data.page,
		albums: res.data.albums,
		filter,
	};
};
