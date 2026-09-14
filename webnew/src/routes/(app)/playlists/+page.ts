import { getPagedQueryOptions } from "$lib/utils";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";
import { FullFilter } from "./types";

function constructFilterSort(
  filter: FullFilter,
  query: Record<string, string>,
  currentUserId?: string,
) {
  const filters = [];

  if (filter.query !== "") {
    filters.push(`name contains "${filter.query}"`);
  }

  if (!filter.filters.all && currentUserId) {
    filters.push(`ownerId = "${currentUserId}"`);
  }

  query["filter"] = filters.join(" and ");

  switch (filter.sort) {
    case "name-a-z":
      query["sort"] = "+name";
      break;
    case "name-z-a":
      query["sort"] = "-name";
      break;
    case "tracks-most":
      query["sort"] = "-trackCount";
      break;
    case "tracks-least":
      query["sort"] = "+trackCount";
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
      all: url.searchParams.get("all") ?? undefined,
    },
    excludes: {},
  });

  constructFilterSort(filter, query, data.user?.id);

  console.log(query);

  const res = await data.apiClient.getPlaylists({ query });
  if (!res.success) {
    throw error(res.error.code, { message: res.error.message });
  }

  return {
    ...data,
    page: res.data.page,
    playlists: res.data.playlists,
    filter,
  };
};
