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

	const albums = await data.apiClient.getUserYearReviewMonthAlbums(
		params.id,
		String(year),
		String(month),
	);
	if (!albums.success) {
		throw error(albums.error.code, { message: albums.error.message });
	}

	return {
		...data,
		year,
		month,

		page: albums.data.page,
		albums: albums.data.albums,
	};
};
