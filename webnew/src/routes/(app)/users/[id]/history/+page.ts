import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent }) => {
	const data = await parent();

	const history = await data.apiClient.getTrackHistory({});
	if (!history.success) {
		throw error(history.error.code, { message: history.error.message });
	}

	return {
		...data,
		page: history.data.page,
		history: history.data.history,
	};
};
