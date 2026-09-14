import type { ApiClient } from "$lib/api/client";
import type { PageLoad } from "./$types";

async function getPlaylists(apiClient: ApiClient, userId?: string) {
  const res = await apiClient.getPlaylists({
    query: { filter: `ownerId = "${userId}"` },
  });
  if (!res.success) {
    return [];
  }

  return res.data.playlists;
}

async function getFavorites(apiClient: ApiClient, userId?: string) {
  if (!userId) return [];

  const res = await apiClient.getUserTrackFavoritesById(userId, {
    query: { sort: "-added", perPage: "10" },
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
    playlists: getPlaylists(data.apiClient, data.user?.id),
    favorites: getFavorites(data.apiClient, data.user?.id),
  };
};
