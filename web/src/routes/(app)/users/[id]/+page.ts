import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, params }) => {
  const data = await parent();

  const stats = await data.apiClient.getUserStats(params.id);
  if (!stats.success) {
    throw error(stats.error.code, { message: stats.error.message });
  }

  const topTracks = await data.apiClient.getUserTopTracks(params.id);
  if (!topTracks.success) {
    throw error(topTracks.error.code, { message: topTracks.error.message });
  }

  const yearStats = await data.apiClient.getUserYearStats(params.id);
  if (!yearStats.success) {
    throw error(yearStats.error.code, { message: yearStats.error.message });
  }

  return {
    ...data,
    stats: stats.data,
    topTracks: topTracks.data.tracks,
    yearStats: yearStats.data.stats,
  };
};
