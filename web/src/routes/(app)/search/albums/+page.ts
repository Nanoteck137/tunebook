import type { Album, Page } from "$lib/api/types";
import { getPagedQueryOptions } from "$lib/utils";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, url }) => {
  const data = await parent();

  const query = url.searchParams.get("query") ?? "";
  const paged = getPagedQueryOptions(url.searchParams);

  let albums = [] as Album[];
  let page: Page | null = null;

  if (query) {
    const res = await data.apiClient.searchAlbums({
      query: { query, ...paged },
    });

    if (res.success) {
      albums = res.data.albums;
      page = res.data.page;
    }
  }

  return {
    ...data,
    query,
    albums,
    page,
  };
};
