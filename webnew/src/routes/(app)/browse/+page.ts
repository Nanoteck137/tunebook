import type { ApiClient } from "$lib/api/client";
import type { PageLoad } from "./$types";

async function getRecentAlbums(apiClient: ApiClient) {
  const res = await apiClient.getAlbums({
    query: { sort: "-created", perPage: "12" },
  });
  if (!res.success) {
    return [];
  }

  return res.data.albums;
}

async function getRecentTracks(apiClient: ApiClient) {
  const res = await apiClient.getTracks({
    query: { sort: "-created", perPage: "12" },
  });
  if (!res.success) {
    return [];
  }

  return res.data.tracks;
}

async function getRecentPlaylists(apiClient: ApiClient) {
  const res = await apiClient.getPlaylists({
    query: { sort: "-created", perPage: "12" },
  });
  if (!res.success) {
    return [];
  }

  return res.data.playlists;
}

export const load: PageLoad = async ({ parent }) => {
  const data = await parent();

  const [recentAlbums, recentTracks, recentPlaylists] = await Promise.all([
    getRecentAlbums(data.apiClient),
    getRecentTracks(data.apiClient),
    getRecentPlaylists(data.apiClient),
  ]);

  return {
    ...data,
    recentAlbums,
    recentTracks,
    recentPlaylists,
  };
};