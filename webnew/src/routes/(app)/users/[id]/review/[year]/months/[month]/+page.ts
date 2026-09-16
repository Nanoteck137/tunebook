import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, params }) => {
	const data = await parent();

	const year = Number(params.year);
	if (!Number.isInteger(year)) {
		throw error(400, { message: "Invalid year" });
	}

	const month = Number(params.month);
	if (!Number.isInteger(month) || month < 1 || month > 12) {
		throw error(400, { message: "Invalid month" });
	}

	const review = await data.apiClient.getUserYearReview(params.id, String(year));
	if (!review.success) {
		throw error(review.error.code, { message: review.error.message });
	}

	const monthData = review.data.months[month - 1] ?? null;
	if (!monthData || monthData.playCount === 0) {
		throw error(404, { message: "Month not found" });
	}

	return {
		...data,
		year,
		month,
		review: review.data,
		monthData,
	};
};