import type { Album, Page } from "$lib/api/types";
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
    filters.push(`name % "%${filter.query}%"`);
  }

  if (filter.filters.decade !== "none") {
    switch (filter.filters.decade) {
      case "1960":
        filters.push("year >= 1960 && year <= 1969");
        break;
      case "1970":
        filters.push("year >= 1970 && year <= 1979");
        break;
      case "1980":
        filters.push("year >= 1980 && year <= 1989");
        break;
      case "1990":
        filters.push("year >= 1990 && year <= 1999");
        break;
      case "2000":
        filters.push("year >= 2000 && year <= 2009");
        break;
      case "2010":
        filters.push("year >= 2010 && year <= 2019");
        break;
      case "2020":
        filters.push("year >= 2020 && year <= 2029");
        break;
    }
  }

  if (filter.filters.tags.length > 0) {
    const s = filter.filters.tags.map((i) => `"${i}"`).join(",");
    filters.push(`hasTag(${s})`);
  }

  if (filter.excludes.tags.length > 0) {
    const s = filter.excludes.tags.map((i) => `"${i}"`).join(",");
    filters.push(`!hasTag(${s})`);
  }

  query["filter"] = filters.join(" && ");

  switch (filter.sort) {
    case "name-a-z":
      query["sort"] = "sort=+name";
      break;
    case "name-z-a":
      query["sort"] = "sort=-name";
      break;
    case "created-new":
      query["sort"] = "sort=-created";
      break;
    case "created-old":
      query["sort"] = "sort=+created";
      break;
    case "updated-new":
      query["sort"] = "sort=-updated";
      break;
    case "updated-old":
      query["sort"] = "sort=+updated";
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
      decade: url.searchParams.get("decade") ?? undefined,
      tags: url.searchParams.get("tags")?.split(",") ?? [],
    },
    excludes: {
      tags: url.searchParams.get("excludeTags")?.split(",") ?? [],
    },
  });

  // console.log(filter);

  constructFilterSort(filter, query);

  // console.log(query);

  let albums: Album[] = [];
  let page: Page | null = null;

  const res = await data.apiClient.getAlbums({ query });
  if (!res.success) {
    throw error(res.error.code, { message: res.error.message });
  }

  albums = res.data.albums;
  page = res.data.page;

  return {
    ...data,
    page,
    albums,
    filter,
  };
};
