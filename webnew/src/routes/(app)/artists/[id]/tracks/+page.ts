import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";
import { TrackFilter, applySort } from "./types";

export const load: PageLoad = async ({ parent, params, url }) => {
	const data = await parent();

	const query: Record<string, string> = {};

	const filter = TrackFilter.parse({
		sort: url.searchParams.get("sort") ?? undefined,
	});
	applySort(filter, query);

	const tracks = await data.apiClient.getTracks({
		query: {
			...query,
			filter: `artistId = "${params.id}" or featuringArtists has "${params.id}"`,
		},
	});
	if (!tracks.success) {
		throw error(tracks.error.code, {
			message: tracks.error.message,
		});
	}

	return {
		...data,
		page: tracks.data.page,
		tracks: tracks.data.tracks,
		filter,
	};
};
