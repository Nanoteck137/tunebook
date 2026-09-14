import { error } from "@sveltejs/kit";
import type { LayoutLoad } from "./$types";

export const load: LayoutLoad = async ({ parent, params }) => {
  const data = await parent();

  const res = await data.apiClient.getArtistById(params.id);
  if (!res.success) {
    throw error(res.error.code, { message: res.error.message });
  }

  return {
    ...data,
    artist: res.data.artist,
  };
};
