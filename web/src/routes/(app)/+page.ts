import type { ApiClient } from "$lib/api/client";
import type { PageLoad } from "./$types";
import type { GetUserStats, Playlist, Track } from "$lib/api/types";

async function getPlaylists(apiClient: ApiClient, userId?: string) {
  if (!userId) return [];

  const res = await apiClient.getPlaylists({
    query: { filter: `ownerId = "${userId}"`, sort: "position" },
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

async function getRecentlyPlayed(apiClient: ApiClient): Promise<Track[]> {
  const res = await apiClient.getTrackHistory({
    query: { perPage: "20", sort: "-listenedAt" },
  });
  if (!res.success) {
    return [];
  }

  const seen = new Set<string>();
  const tracks: Track[] = [];
  for (const entry of res.data.history) {
    if (seen.has(entry.track.id)) continue;
    seen.add(entry.track.id);
    tracks.push(entry.track);
  }

  return tracks.slice(0, 10);
}

async function getTopTracks(
  apiClient: ApiClient,
  userId?: string,
): Promise<Track[]> {
  if (!userId) return [];

  const res = await apiClient.getUserTotalReviewTracks(userId, {
    query: { perPage: "10" },
  });
  if (!res.success) {
    return [];
  }

  return res.data.tracks;
}

async function getStats(
  apiClient: ApiClient,
  userId?: string,
): Promise<GetUserStats | null> {
  if (!userId) return null;

  const res = await apiClient.getUserStats(userId);
  if (!res.success) {
    return null;
  }

  return res.data;
}

export const load: PageLoad = async ({ parent }) => {
  const data = await parent();

  return {
    ...data,
    playlists: getPlaylists(data.apiClient, data.user?.id),
    favorites: getFavorites(data.apiClient, data.user?.id),
    recentlyPlayed: getRecentlyPlayed(data.apiClient),
    topTracks: getTopTracks(data.apiClient, data.user?.id),
    stats: getStats(data.apiClient, data.user?.id),
  };
};
