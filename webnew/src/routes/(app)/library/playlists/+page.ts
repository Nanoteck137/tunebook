import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";
import { FullFilter, constructFilterSort } from "../../playlists/types";

export const load: PageLoad = async ({ parent, url }) => {
	const data = await parent();

	const query: Record<string, string> = {};

	const filter = FullFilter.parse({
		query: url.searchParams.get("query") ?? "",
		sort: url.searchParams.get("sort") ?? undefined,
		filters: {},
		excludes: {},
	});

	constructFilterSort(filter, query, data.user?.id);

	const res = await data.apiClient.getPlaylists({ query });
	if (!res.success) {
		throw error(res.error.code, { message: res.error.message });
	}

	return {
		...data,
		page: res.data.page,
		playlists: res.data.playlists,
		filter,
	};
};
