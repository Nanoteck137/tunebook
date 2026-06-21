import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, params }) => {
  const data = await parent();

  const stats = await data.apiClient.getUserStats(params.id);
  if (!stats.success) {
    throw error(stats.error.code, { message: stats.error.message });
  }

  return {
    ...data,
    stats: stats.data,
  };
};
