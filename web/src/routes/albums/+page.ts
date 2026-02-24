import { getPagedQueryOptions } from "$lib/utils";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, url }) => {
  const data = await parent();

  // TODO(patrik): Fix this
  const query = getPagedQueryOptions(url.searchParams);
  const albums = await data.apiClient.getAlbums({
    query,
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
