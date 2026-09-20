import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, params }) => {
	const data = await parent();

	const res = await data.apiClient.getPlaylists({
		query: {
			filter: `ownerId = "${params.id}"`,
			sort: "position",
		},
	});
	if (!res.success) {
		throw error(res.error.code, { message: res.error.message });
	}

	return {
		...data,
		page: res.data.page,
		playlists: res.data.playlists,
	};
};
