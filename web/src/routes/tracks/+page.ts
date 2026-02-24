import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, url }) => {
  const data = await parent();

  const searchParams = url.searchParams;

  const query: Record<string, string> = {};

  const page = searchParams.get("page");
  if (page) {
    query["page"] = page;
  }

  const tracks = await data.apiClient.getTracks({ query });
  if (!tracks.success) {
    throw error(tracks.error.code, { message: tracks.error.message });
  }

  return {
    ...data,
    page: tracks.data.page,
    tracks: tracks.data.tracks,
  };
};
