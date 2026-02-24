import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent }) => {
  const data = await parent();

  const taglists = await data.apiClient.getTaglists();
  if (!taglists.success) {
    throw error(taglists.error.code, { message: taglists.error.message });
  }

  return {
    ...data,
    taglists: taglists.data.taglists,
  };
};
