import { getPagedQueryOptions } from "$lib/utils";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, params, url }) => {
  const data = await parent();

  const playlist = await data.apiClient.getPlaylistById(params.id);
  if (!playlist.success) {
    throw error(playlist.error.code, { message: playlist.error.message });
  }

  // TODO(patrik): Change getPagedQueryOptions?
  const query = getPagedQueryOptions(url.searchParams);
  const items = await data.apiClient.getPlaylistItems(params.id, {
    query,
  });
  if (!items.success) {
    throw error(items.error.code, { message: items.error.message });
  }

  return {
    ...data,
    playlist: playlist.data,
    page: items.data.page,
    items: items.data.items,
  };
};
