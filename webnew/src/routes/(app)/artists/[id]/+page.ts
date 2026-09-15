import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, params }) => {
	const data = await parent();
	const id = params.id;

	const [albums, tracks, featuredAlbums, featuredTracks] =
		await Promise.all([
			data.apiClient.getAlbums({
				query: { filter: `artistId = "${id}"`, perPage: "6" },
			}),
			data.apiClient.getTracks({
				query: { filter: `artistId = "${id}"`, perPage: "5" },
			}),
			data.apiClient.getAlbums({
				query: { filter: `featuringArtists has "${id}"`, perPage: "6" },
			}),
			data.apiClient.getTracks({
				query: { filter: `featuringArtists has "${id}"`, perPage: "5" },
			}),
		]);

	if (!albums.success) {
		throw error(albums.error.code, { message: albums.error.message });
	}

	if (!tracks.success) {
		throw error(tracks.error.code, { message: tracks.error.message });
	}

	if (!featuredAlbums.success) {
		throw error(featuredAlbums.error.code, {
			message: featuredAlbums.error.message,
		});
	}

	if (!featuredTracks.success) {
		throw error(featuredTracks.error.code, {
			message: featuredTracks.error.message,
		});
	}

	return {
		...data,
		albums: albums.data.albums,
		albumPage: albums.data.page,
		tracks: tracks.data.tracks,
		trackPage: tracks.data.page,
		featuredAlbums: featuredAlbums.data.albums,
		featuredAlbumPage: featuredAlbums.data.page,
		featuredTracks: featuredTracks.data.tracks,
		featuredTrackPage: featuredTracks.data.page,
	};
};