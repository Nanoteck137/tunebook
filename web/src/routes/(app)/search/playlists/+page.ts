import type { Page, Playlist } from "$lib/api/types";
import { getPagedQueryOptions } from "$lib/utils";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, url }) => {
  const data = await parent();

  const query = url.searchParams.get("query") ?? "";
  const paged = getPagedQueryOptions(url.searchParams);

  let playlists = [] as Playlist[];
  let page: Page | null = null;

  if (query) {
    const res = await data.apiClient.searchPlaylists({
      query: { query, ...paged },
    });

    if (res.success) {
      playlists = res.data.playlists;
      page = res.data.page;
    }
  }

  return {
    ...data,
    query,
    playlists,
    page,
  };
};
