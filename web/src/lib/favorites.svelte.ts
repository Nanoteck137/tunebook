import { handleApiError } from "$lib";
import type { ApiClient } from "$lib/api/client";
import { getContext, setContext } from "svelte";

class Favorites {
  apiClient: ApiClient;
  ids = $state<string[]>([]);

  constructor(apiClient: ApiClient, ids: string[]) {
    this.apiClient = apiClient;
    this.ids = ids;
  }

  async toggleTrack(trackId: string) {
    if (this.hasTrack(trackId)) {
      const res = await this.apiClient.unfavoriteTrack(trackId);

      if (!res.success) {
        handleApiError(res.error);
        return;
      }
    } else {
      const res = await this.apiClient.favoriteTrack(trackId);

      if (!res.success) {
        handleApiError(res.error);
        return;
      }
    }

    const ids = await this.apiClient.getFavoriteTrackIds();
    if (!ids.success) {
      handleApiError(ids.error);
      return;
    }

    this.ids = ids.data.ids;
  }

  hasTrack(trackId: string) {
    return !!this.ids.find((v) => v === trackId);
  }
}

const FAVORITES_KEY = Symbol("FAVORITES");

export function setFavorites(apiClient: ApiClient, ids: string[]) {
  return setContext(FAVORITES_KEY, new Favorites(apiClient, ids));
}

export function getFavorites() {
  return getContext<ReturnType<typeof setFavorites>>(FAVORITES_KEY);
}
