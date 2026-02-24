import { getPagedQueryOptions } from "$lib/utils";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, params, url }) => {
  const data = await parent();

  const query = getPagedQueryOptions(url.searchParams);

  const albums = await data.apiClient.getAlbums({
    query: {
      ...query,
      filter: `artistId == "${params.id}" || hasFeaturingArtist("${params.id}")`,
    },
  });
  if (!albums.success) {
    throw error(albums.error.code, { message: albums.error.message });
  }

  return {
    ...data,
    page: albums.data.page,
    albums: albums.data.albums,
  };
};
