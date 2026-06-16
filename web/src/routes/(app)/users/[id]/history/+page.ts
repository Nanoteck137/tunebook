import { getPagedQueryOptions } from "$lib/utils";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, url }) => {
  const data = await parent();

  const query = getPagedQueryOptions(url.searchParams);

  const history = await data.apiClient.getTrackHistory({ query });
  if (!history.success) {
    error(history.error.code, { message: history.error.message });
  }

  return {
    ...data,
    page: history.data.page,
    history: history.data.history,
  };
};
