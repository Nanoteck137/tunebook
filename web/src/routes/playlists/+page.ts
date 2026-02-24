import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent }) => {
  const data = await parent();

  const playlists = await data.apiClient.getPlaylists();
  if (!playlists.success) {
    throw error(playlists.error.code, { message: playlists.error.message });
  }

  return {
    ...data,
    playlists: playlists.data.playlists,
  };
};
