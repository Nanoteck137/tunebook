import { handleApiError } from "$lib";
import type { ApiClient } from "$lib/api/client";
import type { Playlist } from "$lib/api/types";
import { getContext, setContext } from "svelte";
import toast from "svelte-5-french-toast";

class QuickPlaylist {
  apiClient: ApiClient;
  playlist = $state<Playlist | null>(null);
  ids = $state<string[]>([]);

  constructor(apiClient: ApiClient) {
    this.apiClient = apiClient;
  }

  async setPlaylistId(playlistId: string | null) {
    if (playlistId === null) {
      this.playlist = null;
      return;
    }

    const res = await this.apiClient.getPlaylistById(playlistId);
    if (!res.success) {
      // TODO(patrik): Handle error
      handleApiError(res.error);
      return;
    }

    this.playlist = res.data.playlist;

    await this.fetchIds();
  }

  async fetchIds() {
    if (this.playlist === null) {
      this.ids = [];
      return;
    }

    const ids = await this.apiClient.getPlaylistItemIds(this.playlist.id);
    if (!ids.success) {
      handleApiError(ids.error);
      return;
    }

    this.ids = ids.data.ids;
  }

  async toggleTrack(trackId: string) {
    if (this.playlist === null) {
      toast.error("No playlist set for quick playlist");
      return;
    }

    if (this.hasTrack(trackId)) {
      const res = await this.apiClient.removePlaylistItem(this.playlist.id, {
        trackId: trackId,
      });

      if (!res.success) {
        handleApiError(res.error);
        return;
      }
    } else {
      const res = await this.apiClient.addItemToPlaylist(this.playlist.id, {
        trackId: trackId,
      });

      if (!res.success) {
        handleApiError(res.error);
        return;
      }
    }

    await this.fetchIds();
  }

  hasTrack(trackId: string) {
    return !!this.ids.find((v) => v === trackId);
  }
}

const QUICK_PLAYLIST_KEY = Symbol("QUICK_PLAYLIST");

export function setQuickPlaylist(apiClient: ApiClient) {
  return setContext(QUICK_PLAYLIST_KEY, new QuickPlaylist(apiClient));
}

export function getQuickPlaylist() {
  return getContext<ReturnType<typeof setQuickPlaylist>>(QUICK_PLAYLIST_KEY);
}
