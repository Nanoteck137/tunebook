import { getPagedQueryOptions } from "$lib/utils";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

const sortTypes = [
  "name-a-z",
  "name-z-a",
  "year-new",
  "year-old",
  "created-new",
  "created-old",
  "updated-new",
  "updated-old",
] as const;
type SortType = (typeof sortTypes)[number];

function applySort(sort: SortType, query: Record<string, string>) {
  switch (sort) {
    case "name-a-z":
      query["sort"] = "+name";
      break;
    case "name-z-a":
      query["sort"] = "-name";
      break;
    case "year-new":
      query["sort"] = "-year";
      break;
    case "year-old":
      query["sort"] = "+year";
      break;
    case "created-new":
      query["sort"] = "-created";
      break;
    case "created-old":
      query["sort"] = "+created";
      break;
    case "updated-new":
      query["sort"] = "-updated";
      break;
    case "updated-old":
      query["sort"] = "+updated";
      break;
  }
}

export const load: PageLoad = async ({ parent, params, url }) => {
  const data = await parent();

  const query = getPagedQueryOptions(url.searchParams);

  const sort = (url.searchParams.get("sort") ?? "name-a-z") as SortType;
  applySort(sort, query);

  const albums = await data.apiClient.getAlbums({
    query: {
      ...query,
      filter: `artistId = "${params.id}" or featuringArtists has "${params.id}"`,
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
