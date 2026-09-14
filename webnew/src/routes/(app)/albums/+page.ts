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

  if (filter.filters.decade !== "none") {
    switch (filter.filters.decade) {
      case "1960":
        filters.push("(year >= 1960 and year <= 1969)");
        break;
      case "1970":
        filters.push("(year >= 1970 and year <= 1979)");
        break;
      case "1980":
        filters.push("(year >= 1980 and year <= 1989)");
        break;
      case "1990":
        filters.push("(year >= 1990 and year <= 1999)");
        break;
      case "2000":
        filters.push("(year >= 2000 and year <= 2009)");
        break;
      case "2010":
        filters.push("(year >= 2010 and year <= 2019)");
        break;
      case "2020":
        filters.push("(year >= 2020 and year <= 2029)");
        break;
    }
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
      decade: url.searchParams.get("decade") ?? undefined,
      tags: url.searchParams.get("tags")?.split(",") ?? [],
    },
    excludes: {
      tags: url.searchParams.get("excludeTags")?.split(",") ?? [],
    },
  });

  constructFilterSort(filter, query);

  console.log(query);

  const res = await data.apiClient.getAlbums({ query });
  if (!res.success) {
    throw error(res.error.code, { message: res.error.message });
  }

  return {
    ...data,
    page: res.data.page,
    albums: res.data.albums,
    filter,
  };
};
