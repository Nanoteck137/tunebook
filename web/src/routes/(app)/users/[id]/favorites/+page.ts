import { getPagedQueryOptions } from "$lib/utils";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, params, url }) => {
  const data = await parent();

  const query = getPagedQueryOptions(url.searchParams);

  const favorites = await data.apiClient.getUserTrackFavorites(params.id, {
    query,
  });
  if (!favorites.success) {
    throw error(favorites.error.code, { message: favorites.error.message });
  }

  return {
    ...data,
    page: favorites.data.page,
    tracks: favorites.data.items,
  };
};
