import { PUBLIC_API_ADDRESS } from "$env/static/public";
import { setApiClientAuth } from "$lib";
import { ApiClient } from "$lib/api/client";
import type { GetMe, Playlist } from "$lib/api/types";
import { error } from "@sveltejs/kit";
import type { LayoutLoad } from "./$types";

export const prerender = false;
export const ssr = false;

export const load: LayoutLoad = async ({ url }) => {
  console.log("LAYOUT");

  let addr = PUBLIC_API_ADDRESS;
  if (addr === "") {
    addr = url.origin;
  }

  const apiClient = new ApiClient(addr);
  const token = localStorage.getItem("token") ?? undefined;
  setApiClientAuth(apiClient, token);

  let user: GetMe | null = null;
  if (token) {
    const res = await apiClient.getMe();
    if (!res.success) {
      console.error("Get Me API Error", res.error.message);
      user = null;
    } else {
      user = res.data;
    }
  }

  let favoriteIds = [] as string[];
  let quickPlaylistIds = [] as string[];
  let userPlaylists: Playlist[] | null = null;

  if (user) {
    const res = await apiClient.getUserFavoritesIds();
    if (!res.success) {
      // TODO(patrik): Better handling of this error
      throw error(res.error.code, { message: res.error.message });
    }

    favoriteIds = res.data.trackIds;

    if (user.quickPlaylist) {
      const res = await apiClient.getUserQuickPlaylistItemIds();
      if (!res.success) {
        // TODO(patrik): Better handling of this error
        throw error(res.error.code, { message: res.error.message });
      }

      quickPlaylistIds = res.data.trackIds;
    }

    {
      const res = await apiClient.getPlaylists();
      if (!res.success) {
        // TODO(patrik): Better handling of this error
        throw error(res.error.code, { message: res.error.message });
      }

      userPlaylists = res.data.playlists;
    }
  }

  return {
    apiClient,
    user,
    favoriteIds,
    quickPlaylistIds,
    userPlaylists,
  };
};
