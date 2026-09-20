import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, params }) => {
	const data = await parent();

	const review = await data.apiClient.getUserTotalReview(params.id);
	if (!review.success) {
		throw error(review.error.code, { message: review.error.message });
	}

	const topTracks = await data.apiClient.getUserTotalReviewTracks(params.id, {
		query: {
			perPage: "10",
		},
	});
	if (!topTracks.success) {
		throw error(topTracks.error.code, { message: topTracks.error.message });
	}

	const topAlbums = await data.apiClient.getUserTotalReviewAlbums(params.id, {
		query: {
			perPage: "10",
		},
	});
	if (!topAlbums.success) {
		throw error(topAlbums.error.code, { message: topAlbums.error.message });
	}

	const topArtists = await data.apiClient.getUserTotalReviewArtists(
		params.id,
		{
			query: {
				perPage: "10",
			},
		},
	);
	if (!topArtists.success) {
		throw error(topArtists.error.code, { message: topArtists.error.message });
	}

	const topTags = await data.apiClient.getUserTotalReviewTags(params.id, {
		query: {
			perPage: "10",
		},
	});
	if (!topTags.success) {
		throw error(topTags.error.code, { message: topTags.error.message });
	}

	const topDecades = await data.apiClient.getUserTotalReviewDecades(
		params.id,
		{
			query: {
				perPage: "10",
			},
		},
	);
	if (!topDecades.success) {
		throw error(topDecades.error.code, { message: topDecades.error.message });
	}

	const months = await data.apiClient.getUserTotalReviewMonths(params.id);
	if (!months.success) {
		throw error(months.error.code, { message: months.error.message });
	}

	return {
		...data,
		review: review.data.review,

		topTrackCount: topTracks.data.page.totalItems,
		topTracks: topTracks.data.tracks,

		topAlbumCount: topAlbums.data.page.totalItems,
		topAlbums: topAlbums.data.albums,

		topArtistCount: topArtists.data.page.totalItems,
		topArtists: topArtists.data.artists,

		topTagCount: topTags.data.page.totalItems,
		topTags: topTags.data.tags,

		topDecadeCount: topDecades.data.page.totalItems,
		topDecades: topDecades.data.decades,

		months: months.data.months,
	};
};
