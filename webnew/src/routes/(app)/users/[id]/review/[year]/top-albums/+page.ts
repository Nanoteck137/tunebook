import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, params }) => {
	const data = await parent();

	const year = Number(params.year);
	if (!Number.isInteger(year)) {
		throw error(400, { message: "Invalid year" });
	}

	const albums = await data.apiClient.getUserYearReviewAlbums(
		params.id,
		String(year),
	);
	if (!albums.success) {
		throw error(albums.error.code, { message: albums.error.message });
	}

	return {
		...data,
		year,
		page: albums.data.page,
		albums: albums.data.albums,
	};
};
