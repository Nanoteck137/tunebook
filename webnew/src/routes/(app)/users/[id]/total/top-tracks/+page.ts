import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, params }) => {
	const data = await parent();

	const tracks = await data.apiClient.getUserTotalReviewTracks(params.id);
	if (!tracks.success) {
		throw error(tracks.error.code, { message: tracks.error.message });
	}

	return {
		...data,
		page: tracks.data.page,
		tracks: tracks.data.tracks,
	};
};
