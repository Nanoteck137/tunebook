import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent }) => {
	const data = await parent();

	if (!data.user) {
		throw error(401, { message: "Not authenticated" });
	}

	const res = await data.apiClient.getTrackFilters();
	if (!res.success) {
		throw error(res.error.code, { message: res.error.message });
	}

	return {
		...data,
		filters: res.data.filters,
	};
};
