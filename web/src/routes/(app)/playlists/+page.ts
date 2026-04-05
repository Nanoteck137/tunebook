import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent }) => {
  const data = await parent();

  let filter = "";
  if (data.user) {
    filter = `ownerId == "${data.user.id}"`;
  }

  const playlists = await data.apiClient.getPlaylists({ query: { filter } });
  if (!playlists.success) {
    throw error(playlists.error.code, { message: playlists.error.message });
  }

  return {
    ...data,
    playlists: playlists.data.playlists,
  };
};
