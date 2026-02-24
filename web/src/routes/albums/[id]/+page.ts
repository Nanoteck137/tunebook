import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, params }) => {
  const data = await parent();

  const album = await data.apiClient.getAlbumById(params.id);
  if (!album.success) {
    throw error(album.error.code, {
      message: album.error.message,
      type: album.error.type,
    });
  }

  const tracks = await data.apiClient.getAlbumTracks(params.id);
  if (!tracks.success) {
    throw error(tracks.error.code, {
      message: tracks.error.message,
      type: tracks.error.type,
    });
  }

  return {
    ...data,
    album: album.data,
    tracks: tracks.data.tracks,
  };
};
