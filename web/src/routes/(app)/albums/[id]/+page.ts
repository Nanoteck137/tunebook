import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";
import { TrackFilter, buildTrackQuery } from "$lib/components/sort";

export const load: PageLoad = async ({ parent, params, url }) => {
	const data = await parent();

	const album = await data.apiClient.getAlbumById(params.id);
	if (!album.success) {
		throw error(album.error.code, {
			message: album.error.message,
			type: album.error.type,
		});
	}

	const query: Record<string, string> = {};

	const filter = TrackFilter.parse({
		query: url.searchParams.get("query") ?? "",
		sort: url.searchParams.get("sort") ?? undefined,
	});
	buildTrackQuery(filter, "album", query);

	const tracks = await data.apiClient.getAlbumTracks(params.id, { query });
	if (!tracks.success) {
		throw error(tracks.error.code, {
			message: tracks.error.message,
			type: tracks.error.type,
		});
	}

	return {
		...data,
		album: album.data.album,
		tracks: tracks.data.tracks,
		filter,
	};
};
