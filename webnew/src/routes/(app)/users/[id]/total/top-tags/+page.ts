import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, params }) => {
	const data = await parent();

	const tags = await data.apiClient.getUserTotalReviewTags(params.id);
	if (!tags.success) {
		throw error(tags.error.code, { message: tags.error.message });
	}

	return {
		...data,
		page: tags.data.page,
		tags: tags.data.tags,
	};
};
