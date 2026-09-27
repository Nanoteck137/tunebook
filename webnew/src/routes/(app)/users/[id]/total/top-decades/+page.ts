import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, params }) => {
	const data = await parent();

	const decades = await data.apiClient.getUserTotalReviewDecades(params.id);
	if (!decades.success) {
		throw error(decades.error.code, { message: decades.error.message });
	}

	return {
		...data,
		page: decades.data.page,
		decades: decades.data.decades,
	};
};
