import { getPagedQueryOptions } from "$lib/utils";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, params, url }) => {
  const data = await parent();

  const taglist = await data.apiClient.getTaglistById(params.id);
  if (!taglist.success) {
    throw error(taglist.error.code, { message: taglist.error.message });
  }

  // TODO(patrik): Fix this
  const query = getPagedQueryOptions(url.searchParams);
  query["filter"] = taglist.data.filter;

  const tracks = await data.apiClient.getTracks({
    query,
  });
  if (!tracks.success) {
    throw error(tracks.error.code, tracks.error.message);
  }

  return {
    ...data,
    taglist: taglist.data,
    page: tracks.data.page,
    tracks: tracks.data.tracks,
  };
};
