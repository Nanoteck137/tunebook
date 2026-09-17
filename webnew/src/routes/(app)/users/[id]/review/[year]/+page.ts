import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, params }) => {
	const data = await parent();

	const year = Number(params.year);
	if (!Number.isInteger(year)) {
		throw error(400, { message: "Invalid year" });
	}

	const review = await data.apiClient.getUserYearReview(params.id, String(year));
	if (!review.success) {
		throw error(review.error.code, { message: review.error.message });
	}

	console.log("Review", review)

	const topTracks = await data.apiClient.getUserYearReviewTracks(params.id, String(year), {
		query: {
			"perPage": "10",
		}
	});
	if (!topTracks.success) {
		throw error(topTracks.error.code, { message: topTracks.error.message });
	}

	console.log("tracks", topTracks)

	const topAlbums = await data.apiClient.getUserYearReviewAlbums(params.id, String(year), {
		query: {
			"perPage": "10",
		}
	});
	if (!topAlbums.success) {
		throw error(topAlbums.error.code, { message: topAlbums.error.message });
	}

	console.log("albums", topAlbums)

	const topArtists = await data.apiClient.getUserYearReviewArtists(params.id, String(year), {
		query: {
			"perPage": "10",
		}
	});
	if (!topArtists.success) {
		throw error(topArtists.error.code, { message: topArtists.error.message });
	}

	console.log("artists", topArtists)

	const topTags = await data.apiClient.getUserYearReviewTags(params.id, String(year), {
		query: {
			"perPage": "10",
		}
	});
	if (!topTags.success) {
		throw error(topTags.error.code, { message: topTags.error.message });
	}

	console.log("tags", topTags)

	const topDecades = await data.apiClient.getUserYearReviewDecades(params.id, String(year), {
		query: {
			"perPage": "10",
		}
	});
	if (!topDecades.success) {
		throw error(topDecades.error.code, { message: topDecades.error.message });
	}

	console.log("decades", topDecades)

	return {
		...data,
		year,

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
	};
};
