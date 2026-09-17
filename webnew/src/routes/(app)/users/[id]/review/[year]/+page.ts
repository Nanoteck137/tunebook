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

	return {
		...data,
		year,

		review: review.data.review,

		topTrackCount: topTracks.data.page.totalItems,
		topTracks: topTracks.data.tracks,
	};
};
