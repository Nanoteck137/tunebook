import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, params }) => {
	const data = await parent();

	const year = Number(params.year);
	if (!Number.isInteger(year)) {
		throw error(400, { message: "Invalid year" });
	}

	const artists = await data.apiClient.getUserYearReviewArtists(
		params.id,
		String(year),
	);
	if (!artists.success) {
		throw error(artists.error.code, { message: artists.error.message });
	}

	return {
		...data,
		year,
		page: artists.data.page,
		artists: artists.data.artists,
	};
};
