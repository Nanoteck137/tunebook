import type { ApiClient } from "$lib/api/client";
import type { PageLoad } from "./$types";

async function getPlaylists(apiClient: ApiClient) {
  const res = await apiClient.getPlaylists();
  if (!res.success) {
    return [];
  }

  return res.data.playlists;
}

async function getRecentAlbums(apiClient: ApiClient) {
  const res = await apiClient.getAlbums({
    query: { sort: "sort=-created", perPage: "10" },
  });
  if (!res.success) {
    return [];
  }

  return res.data.albums;
}

async function getFavorites(apiClient: ApiClient, userId?: string) {
  if (!userId) return [];

  const res = await apiClient.getUserTrackFavorites(userId, {
    query: { sort: "sort=-added", perPage: "10" },
  });
  if (!res.success) {
    return [];
  }

  return res.data.items;
}

export const load: PageLoad = async ({ parent }) => {
  const data = await parent();

  return {
    ...data,
    playlists: getPlaylists(data.apiClient),
    recentAlbums: getRecentAlbums(data.apiClient),
    favorites: getFavorites(data.apiClient, data.user?.id),
  };
};
