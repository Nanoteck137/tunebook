import { getPagedQueryOptions } from "$lib/utils";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, params, url }) => {
  const data = await parent();

  const query = getPagedQueryOptions(url.searchParams);

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
