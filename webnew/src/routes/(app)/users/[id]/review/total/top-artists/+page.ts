import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, params }) => {
	const data = await parent();

	const artists = await data.apiClient.getUserTotalReviewArtists(params.id);
	if (!artists.success) {
		throw error(artists.error.code, { message: artists.error.message });
	}

	return {
		...data,
		page: artists.data.page,
		artists: artists.data.artists,
	};
};
