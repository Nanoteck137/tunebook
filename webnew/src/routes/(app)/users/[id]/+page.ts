import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, params }) => {
	const data = await parent();

	const stats = await data.apiClient.getUserStats(params.id);
	if (!stats.success) {
		throw error(stats.error.code, { message: stats.error.message });
	}

	const topTracks = await data.apiClient.getUserTotalReviewTracks(params.id, {
		query: {
			perPage: "10",
		},
	});
	if (!topTracks.success) {
		throw error(topTracks.error.code, { message: topTracks.error.message });
	}

	const playlists = await data.apiClient.getPlaylists({
		query: {
			filter: `ownerId = "${params.id}"`,
			sort: "position",
			perPage: "12",
		},
	});
	if (!playlists.success) {
		throw error(playlists.error.code, { message: playlists.error.message });
	}

	const reviews = await data.apiClient.getAllUserYearReviews(params.id);
	if (!reviews.success) {
		throw error(reviews.error.code, { message: reviews.error.message });
	}

	return {
		...data,
		stats: stats.data,
		topTracks: topTracks.data.tracks,
		playlists: playlists.data.playlists,
		reviews: reviews.data.reviews,
	};
};
