import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, params }) => {
	const data = await parent();

	const yearStats = await data.apiClient.getUserYearStats(params.id);
	if (!yearStats.success) {
		throw error(yearStats.error.code, { message: yearStats.error.message });
	}

	return {
		...data,
		yearStats: yearStats.data.stats,
	};
};