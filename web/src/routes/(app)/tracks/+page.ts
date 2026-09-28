import type { TrackFilter } from "$lib/api/types";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";
import { constructFilterSort, FullFilter } from "./types";

export const load: PageLoad = async ({ parent, url }) => {
	const data = await parent();

	let filters: TrackFilter[] | null = null;
	if (data.user) {
		const res = await data.apiClient.getTrackFilters();
		if (!res.success) {
			throw error(res.error.code, { message: res.error.message });
		}

		filters = res.data.filters;
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

	const tracks = await data.apiClient.getTracks({ query });
	if (!tracks.success) {
		throw error(tracks.error.code, { message: tracks.error.message });
	}

	return {
		...data,
		filters,
		page: tracks.data.page,
		tracks: tracks.data.tracks,
		filter,
	};
};
