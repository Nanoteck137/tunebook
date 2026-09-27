import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, params }) => {
	const data = await parent();

	const albums = await data.apiClient.getUserTotalReviewAlbums(params.id);
	if (!albums.success) {
		throw error(albums.error.code, { message: albums.error.message });
	}

	return {
		...data,
		page: albums.data.page,
		albums: albums.data.albums,
	};
};
