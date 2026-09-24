import { getPagedQueryOptions } from "$lib/utils";
import { TrackFilter, buildTrackQuery } from "$lib/components/sort";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, params, url }) => {
	const data = await parent();

	const playlist = await data.apiClient.getPlaylistById(params.id);
	if (!playlist.success) {
		throw error(playlist.error.code, { message: playlist.error.message });
	}

	const query = getPagedQueryOptions(url.searchParams);
	delete query.page;

	const filter = TrackFilter.parse({
		query: url.searchParams.get("query") ?? "",
		sort: url.searchParams.get("sort") ?? undefined,
	});
	buildTrackQuery(filter, "playlist", query);

	const filterId = url.searchParams.get("filterId");
	if (filterId) {
		query["filterId"] = filterId;
	}

	const items = await data.apiClient.getPlaylistItems(params.id, {
		query,
	});
	if (!items.success) {
		throw error(items.error.code, { message: items.error.message });
	}

	return {
		...data,
		playlist: playlist.data.playlist,
		page: items.data.page,
		items: items.data.items,
		filter,
	};
};
