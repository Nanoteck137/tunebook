import { getPagedQueryOptions } from "$lib/utils";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, params, url }) => {
  const data = await parent();

  const virtualPlaylist = await data.apiClient.getVirtualPlaylistById(
    params.id,
  );
  if (!virtualPlaylist.success) {
    throw error(virtualPlaylist.error.code, {
      message: virtualPlaylist.error.message,
    });
  }

  // TODO(patrik): Fix this
  const query = getPagedQueryOptions(url.searchParams);
  const tracks = await data.apiClient.getVirtualPlaylistTracks(params.id, {
    query,
  });
  if (!tracks.success) {
    throw error(tracks.error.code, tracks.error.message);
  }

  return {
    ...data,
    virtualPlaylist: virtualPlaylist.data,
    page: tracks.data.page,
    tracks: tracks.data.tracks,
  };
};
