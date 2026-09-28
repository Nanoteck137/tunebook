import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, params }) => {
	const data = await parent();

	const year = Number(params.year);
	if (!Number.isInteger(year)) {
		throw error(400, { message: "Invalid year" });
	}

	const tracks = await data.apiClient.getUserYearReviewTracks(
		params.id,
		String(year),
	);
	if (!tracks.success) {
		throw error(tracks.error.code, { message: tracks.error.message });
	}

	return {
		...data,
		year,
		page: tracks.data.page,
		tracks: tracks.data.tracks,
	};
};
