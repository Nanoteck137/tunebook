import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, params }) => {
	const data = await parent();

	const reviews = await data.apiClient.getAllUserYearReviews(params.id);
	if (!reviews.success) {
		throw error(reviews.error.code, { message: reviews.error.message });
	}

	// Deliberately not thrown on failure: the year list is the point of this
	// page, so a failed total degrades the All time card instead of blanking
	// the whole index.
	const total = await data.apiClient.getUserTotalReview(params.id);

	return {
		...data,
		reviews: reviews.data.reviews,
		totalReview: total.success ? total.data.review : null,
	};
};
