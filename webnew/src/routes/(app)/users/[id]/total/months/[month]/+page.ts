import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, params }) => {
	const data = await parent();

	const monthNum = Number(params.month);
	if (!Number.isInteger(monthNum) || monthNum < 1 || monthNum > 12) {
		throw error(400, { message: "Invalid month" });
	}

	const months = await data.apiClient.getUserTotalReviewMonths(params.id);
	if (!months.success) {
		throw error(months.error.code, { message: months.error.message });
	}

	const month = months.data.months.find((m) => m.month === monthNum) ?? null;

	return {
		...data,
		monthNum,
		month,
	};
};
