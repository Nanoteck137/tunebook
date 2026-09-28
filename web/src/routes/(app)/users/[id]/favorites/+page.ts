import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";
import {
	FullFilter,
	constructFilterSort,
} from "../../../library/favorites/types";

export const load: PageLoad = async ({ parent, params, url }) => {
	const data = await parent();

	const filters = await data.apiClient.getTrackFilters();
	if (!filters.success) {
		throw error(filters.error.code, { message: filters.error.message });
	}

	const query: Record<string, string> = {};

	const filter = FullFilter.parse({
		query: url.searchParams.get("query") ?? "",
		sort: url.searchParams.get("sort") ?? undefined,
	});

	constructFilterSort(filter, query);

	const filterId = url.searchParams.get("filterId");
	if (filterId) {
		query["filterId"] = filterId;
	}

	const favorites = await data.apiClient.getUserTrackFavoritesById(params.id, {
		query,
	});
	if (!favorites.success) {
		throw error(favorites.error.code, { message: favorites.error.message });
	}

	return {
		...data,
		filters: filters.data.filters,
		page: favorites.data.page,
		tracks: favorites.data.items,
		filter,
	};
};
