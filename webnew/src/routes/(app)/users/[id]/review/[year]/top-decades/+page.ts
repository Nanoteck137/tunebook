import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, params }) => {
	const data = await parent();

	const year = Number(params.year);
	if (!Number.isInteger(year)) {
		throw error(400, { message: "Invalid year" });
	}

	const decades = await data.apiClient.getUserYearReviewDecades(
		params.id,
		String(year),
	);
	if (!decades.success) {
		throw error(decades.error.code, { message: decades.error.message });
	}

	return {
		...data,
		year,
		page: decades.data.page,
		decades: decades.data.decades,
	};
};