import type { Page, Track } from "$lib/api/types";
import { getPagedQueryOptions } from "$lib/utils";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, url }) => {
  const data = await parent();

  const query = url.searchParams.get("query") ?? "";
  const paged = getPagedQueryOptions(url.searchParams);

  let tracks = [] as Track[];
  let page: Page | null = null;

  if (query) {
    const res = await data.apiClient.searchTracks({
      query: { query, ...paged },
    });

    if (res.success) {
      tracks = res.data.tracks;
      page = res.data.page;
    }
  }

  return {
    ...data,
    query,
    tracks,
    page,
  };
};
