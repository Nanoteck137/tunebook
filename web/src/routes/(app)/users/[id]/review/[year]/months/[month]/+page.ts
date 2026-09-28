import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, params }) => {
	const data = await parent();

	const year = Number(params.year);
	if (!Number.isInteger(year)) {
		throw error(400, { message: "Invalid year" });
	}

	const monthNum = Number(params.month);
	if (!Number.isInteger(monthNum) || monthNum < 1 || monthNum > 12) {
		throw error(400, { message: "Invalid month" });
	}

	const review = await data.apiClient.getUserYearReview(
		params.id,
		String(year),
	);
	if (!review.success) {
		throw error(review.error.code, { message: review.error.message });
	}

	const month = await data.apiClient.getUserYearReviewMonth(
		params.id,
		String(year),
		String(monthNum),
	);
	if (!month.success) {
		throw error(month.error.code, { message: month.error.message });
	}

	const topTracks = await data.apiClient.getUserYearReviewMonthTracks(
		params.id,
		String(year),
		String(monthNum),
		{
			query: {
				perPage: "10",
			},
		},
	);
	if (!topTracks.success) {
		throw error(topTracks.error.code, { message: topTracks.error.message });
	}

	const topAlbums = await data.apiClient.getUserYearReviewMonthAlbums(
		params.id,
		String(year),
		String(monthNum),
		{
			query: {
				perPage: "10",
			},
		},
	);
	if (!topAlbums.success) {
		throw error(topAlbums.error.code, { message: topAlbums.error.message });
	}

	const topArtists = await data.apiClient.getUserYearReviewMonthArtists(
		params.id,
		String(year),
		String(monthNum),
		{
			query: {
				perPage: "10",
			},
		},
	);
	if (!topArtists.success) {
		throw error(topArtists.error.code, { message: topArtists.error.message });
	}

	const topTags = await data.apiClient.getUserYearReviewMonthTags(
		params.id,
		String(year),
		String(monthNum),
		{
			query: {
				perPage: "10",
			},
		},
	);
	if (!topTags.success) {
		throw error(topTags.error.code, { message: topTags.error.message });
	}

	const topDecades = await data.apiClient.getUserYearReviewMonthDecades(
		params.id,
		String(year),
		String(monthNum),
		{
			query: {
				perPage: "10",
			},
		},
	);
	if (!topDecades.success) {
		throw error(topDecades.error.code, { message: topDecades.error.message });
	}

	return {
		...data,
		year,
		monthNum,
		review: review.data,
		month: month.data.month,

		topTracks: topTracks.data.tracks,
		topTrackCount: topTracks.data.page.totalItems,

		topAlbums: topAlbums.data.albums,
		topAlbumCount: topAlbums.data.page.totalItems,

		topArtists: topArtists.data.artists,
		topArtistCount: topArtists.data.page.totalItems,

		topTags: topTags.data.tags,
		topTagCount: topTags.data.page.totalItems,

		topDecades: topDecades.data.decades,
		topDecadeCount: topDecades.data.page.totalItems,
	};
};
