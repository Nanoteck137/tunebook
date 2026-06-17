import { getPagedQueryOptions } from "$lib/utils";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

const sortTypes = [
  "name-a-z",
  "name-z-a",
  "artist",
  "album",
  "duration",
  "year",
  "created-new",
  "created-old",
] as const;
type SortType = (typeof sortTypes)[number];

function applySort(sort: SortType, query: Record<string, string>) {
  switch (sort) {
    case "name-a-z":
      query["sort"] = "sort=+name";
      break;
    case "name-z-a":
      query["sort"] = "sort=-name";
      break;
    case "artist":
      query["sort"] = "sort=+artist";
      break;
    case "album":
      query["sort"] = "sort=+album";
      break;
    case "duration":
      query["sort"] = "sort=+duration";
      break;
    case "year":
      query["sort"] = "sort=-year";
      break;
    case "created-new":
      query["sort"] = "sort=-created";
      break;
    case "created-old":
      query["sort"] = "sort=+created";
      break;
  }
}

export const load: PageLoad = async ({ parent, params, url }) => {
  const data = await parent();

  const query = getPagedQueryOptions(url.searchParams);

  const sort = (url.searchParams.get("sort") ?? "name-a-z") as SortType;
  applySort(sort, query);

  const tracks = await data.apiClient.getTracks({
    query: {
      ...query,
      filter: `artistId == "${params.id}" || hasFeaturingArtist("${params.id}")`,
    },
  });
  if (!tracks.success) {
    throw error(tracks.error.code, {
      message: tracks.error.message,
    });
  }

  return {
    ...data,
    page: tracks.data.page,
    tracks: tracks.data.tracks,
  };
};
