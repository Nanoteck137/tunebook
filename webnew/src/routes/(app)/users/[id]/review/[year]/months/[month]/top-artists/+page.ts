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

	const artists = await data.apiClient.getUserYearReviewMonthArtists(
		params.id,
		String(year),
		String(month),
	);
	if (!artists.success) {
		throw error(artists.error.code, { message: artists.error.message });
	}

	return {
		...data,
		year,
		month,

		page: artists.data.page,
		artists: artists.data.artists,
	};
};
