import type { Album, Artist, Track } from "$lib/api/types";
import type { PageLoad } from "./$types";

type Err = {
  code: number;
  message: string;
  type: string;
};

export const load: PageLoad = async ({ parent, url }) => {
  const data = await parent();

  const query = url.searchParams.get("query") ?? "";

  let artists = [] as Artist[];
  let artistError: Err | null = null;

  let albums = [] as Album[];
  let albumError: Err | null = null;

  let tracks = [] as Track[];
  let trackError: Err | null = null;

  const [artistQuery, albumQuery, trackQuery] = await Promise.all([
    data.apiClient.searchArtists({ query: { query } }),
    data.apiClient.searchAlbums({ query: { query } }),
    data.apiClient.searchTracks({ query: { query } }),
  ]);

  if (!artistQuery.success) {
    artistError = artistQuery.error;
  } else {
    artists = artistQuery.data.artists;
  }

  if (!albumQuery.success) {
    albumError = albumQuery.error;
  } else {
    albums = albumQuery.data.albums;
  }

  if (!trackQuery.success) {
    trackError = trackQuery.error;
  } else {
    tracks = trackQuery.data.tracks;
  }

  return {
    ...data,

    query,

    artistError,
    artists,

    albumError,
    albums,

    trackError,
    tracks,
  };
};
