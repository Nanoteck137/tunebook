import { getPagedQueryOptions } from "$lib/utils";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";
import { FullFilter } from "./types";

function constructFilterSort(
  filter: FullFilter,
  query: Record<string, string>,
) {
  const filters = [];

  if (filter.query !== "") {
    filters.push(`name contains "${filter.query}"`);
  }

  filter.filters.tags.forEach((t) => {
    filters.push(`tags has "${t}"`);
  });

  filter.excludes.tags.forEach((t) => {
    filters.push(`not tags has "${t}"`);
  });

  query["filter"] = filters.join(" and ");

  switch (filter.sort) {
    case "name-a-z":
      query["sort"] = "+name";
      break;
    case "name-z-a":
      query["sort"] = "-name";
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

export const load: PageLoad = async ({ parent, url }) => {
  const data = await parent();

  const query = getPagedQueryOptions(url.searchParams);

  const filter = FullFilter.parse({
    query: url.searchParams.get("query") ?? "",
    sort: url.searchParams.get("sort") ?? undefined,
    filters: {
      tags: url.searchParams.get("tags")?.split(",") ?? [],
    },
    excludes: {
      tags: url.searchParams.get("excludeTags")?.split(",") ?? [],
    },
  });

  constructFilterSort(filter, query);

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
