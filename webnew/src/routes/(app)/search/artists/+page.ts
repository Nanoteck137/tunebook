import type { Artist, Page } from "$lib/api/types";
import { getPagedQueryOptions } from "$lib/utils";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, url }) => {
  const data = await parent();

  const query = url.searchParams.get("query") ?? "";
  const paged = getPagedQueryOptions(url.searchParams);

  let artists = [] as Artist[];
  let page: Page | null = null;

  if (query) {
    const res = await data.apiClient.searchArtists({
      query: { query, ...paged },
    });

    if (res.success) {
      artists = res.data.artists;
      page = res.data.page;
    }
  }

  return {
    ...data,
    query,
    artists,
    page,
  };
};
