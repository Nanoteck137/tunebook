import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, params }) => {
  const data = await parent();

  const [albums, tracks] = await Promise.all([
    data.apiClient.getAlbums({
      query: {
        filter: `artistId = "${params.id}" or featuringArtists has "${params.id}"`,
        perPage: "6",
      },
    }),
    data.apiClient.getTracks({
      query: {
        filter: `artistId = "${params.id}" or featuringArtists has "${params.id}"`,
        perPage: "5",
      },
    }),
  ]);

  if (!albums.success) {
    throw error(albums.error.code, {
      message: albums.error.message,
    });
  }

  if (!tracks.success) {
    throw error(tracks.error.code, {
      message: tracks.error.message,
    });
  }

  return {
    ...data,
    albums: albums.data.albums,
    trackPage: tracks.data.page,
    tracks: tracks.data.tracks,
  };
};
