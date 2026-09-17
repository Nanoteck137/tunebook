import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, params }) => {
	const data = await parent();

	const reviews = await data.apiClient.getAllUserYearReviews(params.id);
	if (!reviews.success) {
		throw error(reviews.error.code, { message: reviews.error.message });
	}

	return {
		...data,
		reviews: reviews.data.reviews,
	};
};
