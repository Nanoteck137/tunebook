import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent }) => {
  const data = await parent();

  const virtualPlaylists = await data.apiClient.getVirtualPlaylists();
  if (!virtualPlaylists.success) {
    throw error(virtualPlaylists.error.code, {
      message: virtualPlaylists.error.message,
    });
  }

  return {
    ...data,
    virtualPlaylists: virtualPlaylists.data.virtualPlaylists,
  };
};
