import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, params }) => {
	const data = await parent();

	const year = Number(params.year);
	if (!Number.isInteger(year)) {
		throw error(400, { message: "Invalid year" });
	}

	const top = await data.apiClient.getUserYearReviewTopAlbums(
		params.id,
		String(year),
	);
	if (!top.success) {
		throw error(top.error.code, { message: top.error.message });
	}

	return {
		...data,
		year,
		albums: top.data.albums,
	};
};