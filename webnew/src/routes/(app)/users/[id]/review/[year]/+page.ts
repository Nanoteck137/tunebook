import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, params }) => {
	const data = await parent();

	const year = Number(params.year);
	if (!Number.isInteger(year)) {
		throw error(400, { message: "Invalid year" });
	}

	const yearStats = await data.apiClient.getUserYearStats(params.id);
	if (!yearStats.success) {
		throw error(yearStats.error.code, { message: yearStats.error.message });
	}

	const yearStat = yearStats.data.stats.find((s) => s.year === year) ?? null;

	return {
		...data,
		year,
		yearStat,
	};
};