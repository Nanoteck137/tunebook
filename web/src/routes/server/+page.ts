import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent }) => {
  const data = await parent();

  const mediaSettings = await data.apiClient.getMediaSettings();
  if (!mediaSettings.success) {
    throw error(mediaSettings.error.code, {
      message: mediaSettings.error.message,
      type: mediaSettings.error.type,
    });
  }

  return {
    ...data,
    mediaSettings: mediaSettings.data,
  };
};
