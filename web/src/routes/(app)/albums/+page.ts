import { getPagedQueryOptions } from "$lib/utils";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";
import { FullFilter } from "./types";

export const load: PageLoad = async ({ parent, url }) => {
  const data = await parent();

  // TODO(patrik): Fix this
  const query = getPagedQueryOptions(url.searchParams);

  const filter = FullFilter.parse({
    query: url.searchParams.get("query") ?? "",
    sort: url.searchParams.get("sort") ?? undefined,
    filters: {
      // type: url.searchParams.get("filterType")?.split(",") ?? [],
    },
    excludes: {
      // type: url.searchParams.get("excludeType")?.split(",") ?? [],
    },
  });

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

    filter,
  };
};
