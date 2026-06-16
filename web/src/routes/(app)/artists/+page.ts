import { getPagedQueryOptions } from "$lib/utils";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, url }) => {
  const data = await parent();

  // TODO(patrik): Fix this
  const query = getPagedQueryOptions(url.searchParams);
  const artists = await data.apiClient.getArtists({
    query,
  });
  if (!artists.success) {
    throw error(artists.error.code, artists.error.message);
  }

  return {
    ...data,
    page: artists.data.page,
    artists: artists.data.artists,
  };
};
