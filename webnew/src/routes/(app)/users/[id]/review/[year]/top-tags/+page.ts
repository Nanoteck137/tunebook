import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, params }) => {
	const data = await parent();

	const year = Number(params.year);
	if (!Number.isInteger(year)) {
		throw error(400, { message: "Invalid year" });
	}

	const tags = await data.apiClient.getUserYearReviewTags(
		params.id,
		String(year),
	);
	if (!tags.success) {
		throw error(tags.error.code, { message: tags.error.message });
	}

	return {
		...data,
		year,
		page: tags.data.page,
		tags: tags.data.tags,
	};
};