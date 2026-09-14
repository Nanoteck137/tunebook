import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent }) => {
  const data = await parent();

  const res = await data.apiClient.getAlbums({
    query: { sort: "-created", perPage: "10" },
  });
  if (!res.success) {
    throw error(res.error.code, { message: res.error.message });
  }

  return {
    ...data,
    recentAlbums: res.data.albums,
  };
};