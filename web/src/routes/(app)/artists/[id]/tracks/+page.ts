import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";
import { TrackFilter, buildArtistTracksQuery } from "./types";

export const load: PageLoad = async ({ parent, params, url }) => {
	const data = await parent();

	const filter = TrackFilter.parse({
		query: url.searchParams.get("query") ?? "",
		sort: url.searchParams.get("sort") ?? undefined,
	});

	const query: Record<string, string> = {};
	buildArtistTracksQuery(filter, params.id, query);

	const tracks = await data.apiClient.getTracks({ query });
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
