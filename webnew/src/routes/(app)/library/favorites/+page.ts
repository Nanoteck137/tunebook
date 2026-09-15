import type { TrackFilter } from "$lib/api/types";
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

  if (!data.user) {
    throw error(401, { message: "Not authenticated" });
  }

  let filters: TrackFilter[] | null = null;
  const res = await data.apiClient.getTrackFilters();
  if (!res.success) {
    throw error(res.error.code, { message: res.error.message });
  }
  filters = res.data.filters;

  const query = getPagedQueryOptions(url.searchParams);

  const filter = FullFilter.parse({
    query: url.searchParams.get("query") ?? "",
    sort: url.searchParams.get("sort") ?? undefined,
  });

  constructFilterSort(filter, query);

  const filterId = url.searchParams.get("filterId");
  if (filterId) {
    query["filterId"] = filterId;
  }

  const favorites = await data.apiClient.getUserTrackFavoritesById(
    data.user.id,
    { query },
  );
  if (!favorites.success) {
    throw error(favorites.error.code, { message: favorites.error.message });
  }

  return {
    ...data,
    filters,
    page: favorites.data.page,
    tracks: favorites.data.items,
  };
};